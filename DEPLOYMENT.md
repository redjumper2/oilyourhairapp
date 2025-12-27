# Multi-Site Deployment Guide

Deploy multiple static websites to nginx with SSL/TLS via **Cloudflare Tunnel** - perfect for home networks without port forwarding.

## Repository Structure

```
oilyourhairapp/
├── Makefile                          # Main orchestrator
├── pub-site-setup.mk                 # Generic site deployment targets
├── sites.conf                        # Site configuration (domains list)
├── DEPLOYMENT.md                     # This file
├── etc/
│   └── nginx/
│       └── sites-available/
│           ├── oilyourhair.com      # Nginx config for each site
│           └── example.com
└── var/
    └── www/
        ├── oilyourhair.com/
        │   └── html/                 # HTML files for this site
        └── example.com/
            └── html/
```

## Quick Start

### Deploy Existing Site (oilyourhair.com)

```bash
# Deploy nginx config and HTML files
make oilyourhair-deploy

# Check deployment status
make oilyourhair-status

# View logs
make oilyourhair-logs
```

### Add a New Site

1. **Add site to `sites.conf`:**
   ```bash
   echo "MYSITE_SITE=example.com" >> sites.conf
   ```

2. **Create directory structure:**
   ```bash
   mkdir -p etc/nginx/sites-available
   mkdir -p var/www/example.com/html
   ```

3. **Create nginx config:**
   ```bash
   cp etc/nginx/sites-available/oilyourhair.com etc/nginx/sites-available/example.com
   # Edit the file and change domain names
   ```

4. **Add HTML files:**
   ```bash
   cp -r var/www/oilyourhair.com/html/* var/www/example.com/html/
   # Edit your HTML files
   ```

5. **Deploy:**
   ```bash
   make site-deploy SITE=example.com
   ```

6. **Add convenience targets to Makefile:**
   ```makefile
   mysite-deploy:
       @$(MAKE) site-deploy SITE=$(MYSITE_SITE)

   mysite-status:
       @$(MAKE) site-status SITE=$(MYSITE_SITE)
   ```

7. **Set up Cloudflare Tunnel for the new site:**
   ```bash
   make site-tunnel-create SITE=example.com
   make site-tunnel-config SITE=example.com
   make site-tunnel-test SITE=example.com
   make site-tunnel-route SITE=example.com
   ```

## Cloudflare Tunnel Setup

### Why Cloudflare Tunnel?

- ✅ No port forwarding needed
- ✅ Works behind NAT/firewall (perfect for home networks)
- ✅ Free SSL/TLS (Cloudflare handles certificates)
- ✅ DDoS protection included
- ✅ No exposed ports on your home network

### Initial Setup (One-Time)

```bash
# Complete automated setup
make cloudflare-setup
```

This will:
1. Install cloudflared
2. Login to Cloudflare
3. Create tunnel for oilyourhair.com
4. Configure tunnel settings

### Test and Go Live

```bash
# Test tunnel locally
make oilyourhair-tunnel-test
# Visit the preview URL to verify

# Route DNS (this makes site go live!)
make oilyourhair-tunnel-route

# Install as systemd service (auto-start on boot)
make cloudflare-tunnel-service
```

### Check Tunnel Status

```bash
# View service status
make cloudflare-tunnel-status

# View logs
make cloudflare-tunnel-logs
```

## Common Operations

### Deploy Updates

```bash
# HTML only
make oilyourhair-html-deploy

# Nginx config only
make oilyourhair-conf-deploy

# Both
make oilyourhair-deploy
```

### Monitor Site

```bash
# View deployment info
make oilyourhair-status

# View nginx logs
make oilyourhair-logs

# Check nginx config syntax
sudo nginx -t
```

### Manage Tunnel

```bash
# Restart tunnel
sudo systemctl restart cloudflared

# Stop tunnel
sudo systemctl stop cloudflared

# Start tunnel
sudo systemctl start cloudflared
```

## Migration from Cloud VM to Home Network

Migrating from a cloud VM to your home mini PC:

1. **Set up site locally:**
   ```bash
   make oilyourhair-deploy
   ```

2. **Install and configure Cloudflare Tunnel:**
   ```bash
   make cloudflare-setup
   ```

3. **Test tunnel:**
   ```bash
   make oilyourhair-tunnel-test
   ```

4. **Route DNS (switches live traffic):**
   ```bash
   cloudflared tunnel route dns --overwrite-dns oilyourhair.com-tunnel oilyourhair.com
   cloudflared tunnel route dns --overwrite-dns oilyourhair.com-tunnel www.oilyourhair.com
   ```

5. **Install as service:**
   ```bash
   make cloudflare-tunnel-service
   ```

6. **Verify site is live:**
   ```bash
   curl -I https://oilyourhair.com
   ```

7. **Keep cloud VM running for a few hours as backup, then shut down**

## Makefile Target Reference

### Site-Specific (OilYourHair.com)
- `oilyourhair-deploy` - Full deployment (nginx + html)
- `oilyourhair-conf-deploy` - Deploy nginx config only
- `oilyourhair-html-deploy` - Deploy HTML files only
- `oilyourhair-status` - Check deployment status
- `oilyourhair-logs` - View nginx logs
- `oilyourhair-tunnel-create` - Create Cloudflare Tunnel
- `oilyourhair-tunnel-config` - Configure tunnel
- `oilyourhair-tunnel-route` - Route DNS to tunnel
- `oilyourhair-tunnel-test` - Test tunnel locally

### Generic (Any Site)
- `make site-deploy SITE=domain.com` - Full deployment
- `make site-conf-deploy SITE=domain.com` - Deploy nginx config
- `make site-html-deploy SITE=domain.com` - Deploy HTML files
- `make site-status SITE=domain.com` - Check status
- `make site-tunnel-create SITE=domain.com` - Create tunnel
- `make site-remove SITE=domain.com` - Remove site config

### Infrastructure
- `nginx-install` - Install nginx
- `cloudflare-install` - Install cloudflared
- `cloudflare-login` - Authenticate with Cloudflare
- `cloudflare-setup` - Complete automated setup
- `cloudflare-tunnel-service` - Install tunnel as systemd service
- `cloudflare-tunnel-status` - Check service status
- `cloudflare-tunnel-logs` - View tunnel logs

## How It Works

### Traffic Flow

```
Visitor Browser
    ↓ (HTTPS - Cloudflare's SSL certificate)
Cloudflare Edge Network
    ↓ (Encrypted tunnel)
Your Mini PC (cloudflared)
    ↓ (HTTP on localhost:80)
Nginx → Your Website
```

### SSL/TLS

- **Visitor to Cloudflare:** HTTPS with Cloudflare's certificate (trusted by all browsers)
- **Cloudflare to Mini PC:** Encrypted through Cloudflare Tunnel
- **Tunnel to Nginx:** HTTP on localhost (safe, not exposed to internet)

**Result:** Fully secure HTTPS website without managing certificates!

### Optional: Cloudflare Origin Certificates

For maximum security (encrypt localhost connection):

1. Get certificate from [Cloudflare Dashboard](https://dash.cloudflare.com)
   - Go to: SSL/TLS → Origin Server → Create Certificate
2. Download certificate and private key
3. Install in nginx config:
   ```nginx
   listen 443 ssl;
   ssl_certificate /path/to/cloudflare-origin.pem;
   ssl_certificate_key /path/to/cloudflare-origin-key.pem;
   ```
4. Update tunnel config to use `https://localhost:443`
5. Enable "Full (strict)" SSL mode in Cloudflare

## Troubleshooting

### Nginx errors
```bash
# Test configuration
sudo nginx -t

# Check logs
make oilyourhair-logs

# Restart nginx
sudo systemctl restart nginx
```

### Tunnel not working
```bash
# Check service status
make cloudflare-tunnel-status

# View logs
make cloudflare-tunnel-logs

# Restart service
sudo systemctl restart cloudflared
```

### Site not accessible
```bash
# Check if tunnel is running
make cloudflare-tunnel-status

# Check nginx is running
sudo systemctl status nginx

# Check DNS routing
cloudflared tunnel route dns list
```

### DNS not updating
- DNS changes can take a few minutes to propagate
- Check DNS with: `dig oilyourhair.com`
- Should show CNAME pointing to Cloudflare Tunnel

## Best Practices

1. **Version control nginx configs** - Always commit changes to `etc/nginx/sites-available/`
2. **Test before deploying** - Use `sudo nginx -t` to validate configs
3. **Monitor logs** - Regularly check `make oilyourhair-logs`
4. **Backup tunnel config** - Keep backup of `~/.cloudflared/` directory
5. **Keep systemd service enabled** - Tunnel auto-starts on reboot
6. **Use site-specific targets** - Easier than remembering generic commands

## Security Notes

- Cloudflare Tunnel means **no ports exposed** on your home network
- All traffic is encrypted end-to-end
- DDoS protection included automatically
- No need to manage SSL certificates
- Admin routes can use HTTP basic auth (see `frontend-passwd` target)
