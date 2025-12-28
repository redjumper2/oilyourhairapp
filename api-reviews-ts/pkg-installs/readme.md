What it does
✅ Installs Node.js, MongoDB, Nginx
✅ Creates /home/youruser/myapp with a simple Express app
✅ Creates a systemd service → auto-starts & restarts app
✅ Configures Nginx → reverse-proxy on port 80 → Node.js
✅ (Optionally) sets up free HTTPS with Let’s Encrypt

chmod +x deploy-myapp.sh
3️⃣ Edit these variables at the top of the script if needed:

DOMAIN="yourdomain.com"      # or leave blank for IP
ENABLE_SSL=true              # true/false
APP_USER="yourusername"      # usually your SSH user
