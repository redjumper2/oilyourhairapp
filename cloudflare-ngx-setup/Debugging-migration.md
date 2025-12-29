# Cloudflare Tunnel Debugging Guide

When you see a **502 Bad Gateway** error with "Host error" from Cloudflare, it means Cloudflare can't reach your origin server (mini PC). Here's how to debug it systematically.

## Quick Diagnosis

The error page shows three components:
- **You (working)** - Your browser/client is fine
- **Cloudflare Edge (working)** - Cloudflare's network is fine
- **Origin Host (error)** - Your mini PC/tunnel has an issue

## Debugging Steps

### 1. Check Cloudflared Service Status

```bash
sudo systemctl status cloudflared --no-pager
```

**Look for:**
- `Active: active (running)` - Service is running
- `Registered tunnel connection` messages - Shows active connections to Cloudflare edge
- Should see 4 connections (connIndex 0-3)

### 2. Test Nginx Locally

```bash
curl -I http://localhost:8080
```

**Expected:** `HTTP/1.1 200 OK`
**If fails:** Nginx isn't running or not on port 8080

Check nginx is running:
```bash
sudo systemctl status nginx
```

Check which ports nginx is listening on:
```bash
# Using ss (recommended)
sudo ss -tlnp | grep nginx

# Or using lsof
sudo lsof -i :8080

# Or using netstat (if available)
sudo netstat -tlnp | grep nginx
```

Should see nginx listening on `0.0.0.0:8080` and `[::]:8080`

**Example lsof output:**
```
COMMAND   PID     USER   FD   TYPE DEVICE SIZE/OFF NODE NAME
nginx    2243     root   11u  IPv4  ...      0t0  TCP *:http-alt (LISTEN)
nginx    2243     root   12u  IPv6  ...      0t0  TCP *:http-alt (LISTEN)
```
(http-alt = port 8080)

### 3. Check Cloudflared Logs

```bash
# Recent errors
sudo journalctl -u cloudflared -n 50 --no-pager | grep -E "(ERR|WRN|error)"

# Last 1 minute
sudo journalctl -u cloudflared --since "1 minute ago" --no-pager

# Watch live logs
sudo journalctl -u cloudflared -f
```

**Ignore these warnings (non-critical):**
- `ICMP proxy feature is disabled` - Not needed for HTTP
- `ping_group_range` errors - Not needed for HTTP

**Critical errors to watch for:**
- `Connection terminated`
- `failed to serve tunnel connection`
- `Cannot determine default origin certificate path`
- Errors mentioning `localhost:8080` or `origin`

### 4. Verify Configuration

Check the config file:
```bash
cat /etc/cloudflared/config.yml
```

**Should show:**
```yaml
tunnel: <TUNNEL_ID>
credentials-file: /etc/cloudflared/<TUNNEL_ID>.json

ingress:
  - hostname: oilyourhair.com
    service: http://localhost:8080
  - hostname: www.oilyourhair.com
    service: http://localhost:8080
  - service: http_status:404
```

Validate the configuration:
```bash
sudo cloudflared tunnel ingress validate
```

**Expected:** `OK`

Test what a specific URL routes to:
```bash
sudo cloudflared tunnel ingress rule https://oilyourhair.com
```

**Expected:**
```
Matched rule #0
  hostname: oilyourhair.com
  service: http://localhost:8080
```

### 5. Check Tunnel Info

```bash
sudo cloudflared tunnel info oilyourhair.com-tunnel
```

**Look for:**
- Tunnel ID matches config file
- Shows active connectors
- Edge locations (e.g., dfw06, dfw13)

### 6. Test Website

```bash
# Check HTTP status
curl -I https://oilyourhair.com 2>&1 | grep HTTP

# Verbose output to see connection details
curl -v https://oilyourhair.com 2>&1 | head -30
```

**If 502:** Problem with tunnel → nginx connection
**If timeout:** DNS or tunnel not connected
**If 200:** Everything working!

### 7. Check DNS Routing

```bash
# Check DNS records
dig oilyourhair.com +short
```

**Should show:** Cloudflare IPs (104.x.x.x, 172.x.x.x)

Verify tunnel DNS routing:
```bash
sudo cloudflared tunnel route dns list
```

### 8. Check Firewall (If Applicable)

```bash
# Check UFW status
sudo ufw status

# Check iptables
sudo iptables -L -n | head -20
```

**Note:** Localhost traffic (127.0.0.1) typically bypasses firewalls, but worth checking if you have strict rules.

### 9. Advanced: Run Cloudflared in Debug Mode

If everything looks correct but still getting 502, run cloudflared manually with debug logging:

```bash
# Stop the service
sudo systemctl stop cloudflared

# Run manually with debug logging
sudo cloudflared --config /etc/cloudflared/config.yml --loglevel debug tunnel run

# In another terminal, test the site
curl -I https://oilyourhair.com

# Watch the debug output for:
# - "HEAD https://oilyourhair.com/" - Shows incoming request
# - "200 OK" - Shows successful proxy to nginx
# - Any errors connecting to localhost:8080
```

Press `Ctrl+C` to stop, then restart the service:
```bash
sudo systemctl start cloudflared
```

**If it works in debug mode but not as a service:**
- Service configuration issue
- Stale service state
- Try: `sudo systemctl restart cloudflared`

## Common Issues & Solutions

### Issue: 502 Error but tunnel shows as connected

**Symptoms:**
- `systemctl status cloudflared` shows running
- 4 registered connections
- But website returns 502

**Solution:**
```bash
# Restart the service to clear stale state
sudo systemctl restart cloudflared

# Wait a few seconds
sleep 5

# Test again
curl -I https://oilyourhair.com
```

### Issue: Nginx not responding on port 8080

**Check nginx config:**
```bash
cat /etc/nginx/sites-available/oilyourhair.com | grep listen
```

Should show:
```
listen 8080;
listen [::]:8080;
```

**Verify nginx is actually listening on port 8080:**
```bash
sudo lsof -i :8080
# or
sudo ss -tlnp | grep 8080
```

**If wrong port, fix and reload:**
```bash
# Edit config
sudo nano /etc/nginx/sites-available/oilyourhair.com

# Test config
sudo nginx -t

# Reload nginx
sudo systemctl reload nginx

# Verify it's listening on 8080
sudo lsof -i :8080
```

### Issue: Config changes not taking effect

**Problem:** Config file changed but service still uses old config

**Solution:**
```bash
# Verify service loads correct config
sudo journalctl -u cloudflared --since "5 minutes ago" | grep -i "config\|Settings"

# Should show: config:/etc/cloudflared/config.yml

# If using wrong config, restart service
sudo systemctl restart cloudflared
```

### Issue: Permission errors with cert.pem

**Error:** `cannot read /etc/cloudflared/cert.pem to load origin certificate`

**Check permissions:**
```bash
ls -la /etc/cloudflared/cert.pem
```

**Should be:** `-rw------- 1 root root`

**If wrong, fix:**
```bash
sudo chmod 600 /etc/cloudflared/cert.pem
sudo chown root:root /etc/cloudflared/cert.pem
```

### Issue: Tunnel not starting after boot

**Check service is enabled:**
```bash
sudo systemctl is-enabled cloudflared
```

**If not enabled:**
```bash
sudo systemctl enable cloudflared
```

## Systematic Debugging Checklist

When troubleshooting 502 errors, check in this order:

- [ ] 1. Cloudflared service running? (`systemctl status cloudflared`)
- [ ] 2. Nginx responding locally? (`curl http://localhost:8080`)
- [ ] 3. Nginx listening on 8080? (`ss -tlnp | grep nginx`)
- [ ] 4. Config file correct? (`cat /etc/cloudflared/config.yml`)
- [ ] 5. Config validates? (`cloudflared tunnel ingress validate`)
- [ ] 6. Tunnel connected? (4 connections in systemctl status)
- [ ] 7. DNS routed to tunnel? (`dig oilyourhair.com`)
- [ ] 8. Recent errors in logs? (`journalctl -u cloudflared -n 50`)
- [ ] 9. Try restart? (`systemctl restart cloudflared`)
- [ ] 10. Works in debug mode? (Run manually with `--loglevel debug`)

## Recovery Commands

If everything is broken, here's the nuclear option:

```bash
# 1. Stop everything
sudo systemctl stop cloudflared
sudo systemctl stop nginx

# 2. Start nginx
sudo systemctl start nginx

# Test nginx locally
curl -I http://localhost:8080

# 3. Start cloudflared
sudo systemctl start cloudflared

# Wait 10 seconds for connections to establish
sleep 10

# 4. Test website
curl -I https://oilyourhair.com

# 5. Check logs if still broken
sudo journalctl -u cloudflared -n 100 --no-pager
```

## Getting Help

If still stuck, collect this information:

```bash
# Save diagnostic info to a file
{
  echo "=== Cloudflared Status ==="
  sudo systemctl status cloudflared --no-pager

  echo -e "\n=== Config File ==="
  cat /etc/cloudflared/config.yml

  echo -e "\n=== Tunnel Info ==="
  sudo cloudflared tunnel info oilyourhair.com-tunnel

  echo -e "\n=== Nginx Listening Ports ==="
  sudo ss -tlnp | grep nginx

  echo -e "\n=== Recent Logs ==="
  sudo journalctl -u cloudflared -n 50 --no-pager

  echo -e "\n=== Test Localhost ==="
  curl -I http://localhost:8080

  echo -e "\n=== Test Website ==="
  curl -I https://oilyourhair.com

} > ~/cloudflare-debug.txt

cat ~/cloudflare-debug.txt
```

Share `cloudflare-debug.txt` when asking for help.

## Useful Make Targets

```bash
# Check deployment status
make oilyourhair-status

# View nginx logs
make oilyourhair-logs

# Check tunnel service
make cloudflare-tunnel-status

# View tunnel logs
make cloudflare-tunnel-logs

# Redeploy nginx config
make oilyourhair-conf-deploy

# Redeploy HTML
make oilyourhair-html-deploy
```
