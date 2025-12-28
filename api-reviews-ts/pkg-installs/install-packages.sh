#!/bin/bash
set -e

### CONFIGURATION ###
APP_NAME="myapp"
APP_USER="$USER"                       # assumes current user owns the app
APP_DIR="/home/$APP_USER/myapp"
NODE_VERSION="20"
DOMAIN="example.com"                              # set to yourdomain.com if you have a domain
ENABLE_SSL=true                       # set to true if you want SSL
APP_PORT=3000

### FUNCTIONS ###
log() { echo -e "\e[32m[+] $1\e[0m"; }

### STEP 1: Update & install basics ###
log "Updating system & installing required packages…"
sudo apt update && sudo apt upgrade -y
sudo apt install -y curl software-properties-common build-essential nginx

### STEP 2: Install Node.js ###
log "Installing Node.js $NODE_VERSION…"
curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION}.x | sudo -E bash -
sudo apt install -y nodejs

### STEP 3: Install MongoDB ###
if ! command -v mongod >/dev/null; then
    log "Installing MongoDB…"
    wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -sc)/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list
    sudo apt update
    sudo apt install -y mongodb-org
    sudo systemctl enable mongod --now
else
    log "MongoDB already installed."
fi

### STEP 4: Set up Node.js app ###
log "Setting up Node.js app in $APP_DIR…"
mkdir -p "$APP_DIR"
cd "$APP_DIR"

if [ ! -f package.json ]; then
    log "Initializing Node.js app…"
    npm init -y
    npm install express mongodb
    cat > server.js <<EOF
const express = require('express');
const { MongoClient } = require('mongodb');
const app = express();
const port = ${APP_PORT};
const uri = 'mongodb://localhost:27017';
const dbName = 'testdb';

async function main() {
  const client = new MongoClient(uri, { useUnifiedTopology: true });
  await client.connect();
  console.log('✅ Connected to MongoDB');
  const db = client.db(dbName);
  const collection = db.collection('test');
  if (await collection.countDocuments() === 0) {
    await collection.insertOne({ message: 'Hello from MongoDB!' });
  }
  app.get('/', (req, res) => res.send('Hello World from Node.js + Express 🚀'));
  app.get('/data', async (req, res) => res.json(await collection.find().toArray()));
  app.listen(port, () => console.log(\`🚀 Server running at http://localhost:\${port}\`));
}
main();
EOF
fi

### STEP 5: Create systemd service ###
log "Creating systemd service…"
SERVICE_FILE="/etc/systemd/system/${APP_NAME}.service"

sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=Node.js App: $APP_NAME
After=network.target

[Service]
WorkingDirectory=$APP_DIR
ExecStart=/usr/bin/node $APP_DIR/server.js
Restart=always
User=$APP_USER
Environment=NODE_ENV=production

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now "$APP_NAME"

### STEP 6: Configure Nginx ###
log "Configuring Nginx reverse proxy…"
NGINX_CONF="/etc/nginx/sites-available/${APP_NAME}"

sudo tee "$NGINX_CONF" > /dev/null <<EOF
server {
    listen 80;
    ${DOMAIN:+server_name $DOMAIN;}

    location / {
        proxy_pass http://127.0.0.1:${APP_PORT};
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_cache_bypass \$http_upgrade;
    }
}
EOF

sudo ln -sf "$NGINX_CONF" /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx

### STEP 7: SSL (optional) ###
if [ "$ENABLE_SSL" = true ] && [ -n "$DOMAIN" ]; then
    log "Setting up SSL with Certbot…"
    sudo apt install -y certbot python3-certbot-nginx
    sudo certbot --nginx -d "$DOMAIN" --non-interactive --agree-tos -m "admin@$DOMAIN"
fi

log "✅ Deployment completed!"
log "Visit: http://${DOMAIN:-your-server-ip}/"
[ "$ENABLE_SSL" = true ] && log "or https://${DOMAIN}/"
