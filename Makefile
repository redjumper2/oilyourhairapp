# Project settings
APP_NAME = user-reviews-app
ENTRY = app.js
PORT = 3000
MONGO_URI = mongodb://localhost:27017
MONGO_DB = myappdb

# Include sites configuration
include sites.conf

# Include generic public site setup targets
include pub-site-setup.mk

# Default target
.PHONY: help
help:
	@echo "📦 ${APP_NAME} Makefile"
	@echo ""
	@echo "Backend API Commands:"
	@echo "  install        Install npm dependencies"
	@echo "  run            Run the server"
	@echo "  dev            Run with nodemon (if installed)"
	@echo "  mongo-start    Start MongoDB (macOS: brew only)"
	@echo "  mongo-stop     Stop MongoDB (macOS: brew only)"
	@echo "  mongo-status   Show MongoDB status"
	@echo "  test-api       Show test curl commands"
	@echo ""
	@echo "Public Site Commands:"
	@echo "  oilyourhair-deploy       Deploy oilyourhair.com (nginx + html)"
	@echo "  oilyourhair-conf-deploy  Deploy nginx config only"
	@echo "  oilyourhair-html-deploy  Deploy HTML files only"
	@echo "  oilyourhair-status       Check deployment status"
	@echo "  oilyourhair-logs         View nginx logs"
	@echo ""
	@echo "Generic Site Commands (use SITE=domain.com):"
	@echo "  make site-deploy SITE=example.com"
	@echo "  make site-status SITE=example.com"
	@echo ""
	@echo "Cloudflare Tunnel Commands:"
	@echo "  cloudflare-setup         Complete Cloudflare Tunnel setup"
	@echo "  cloudflare-tunnel-status Check tunnel service status"
	@echo "  cloudflare-tunnel-logs   View tunnel logs"

#------------------------------------
# GLOBAL INFRASTRUCTURE SETUP
#------------------------------------
nginx-install:
	@$(MAKE) site-nginx-install

#------------------------------------
# ADMIN PASSWORD PROTECTION
#------------------------------------
frontend-passwd:
	sudo apt install -y apache2-utils
	sudo htpasswd -c /etc/nginx/.htpasswd oilyourhairadmin
	# Enter password when prompted
	sudo chown www-data:www-data /etc/nginx/.htpasswd

	# to secure route, add:
	# auth_basic "Restricted Access";
	# auth_basic_user_file /etc/nginx/.htpasswd;
	sudo nginx -t
	sudo systemctl reload nginx


#------------------------------------
# use for typescript development
install-dependencies:
	@echo "📦 Installing dependencies..."
	sudo apt purge nodejs npm -y

	curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
	sudo apt-get install -y nodejs
	npm install
# 	@echo "📦 Installing nodemon globally (if not already installed)...
# 	npm install -g nodemon || echo "Nodemon is already installed or failed to install.
	@echo "📦 Installation complete!"

run:
	@echo "🚀 Starting server on port ${PORT}..."
	node $(ENTRY)

deploy:
	@echo "🚀 Deploying ${APP_NAME} to local server..."
	@echo "Running on port ${PORT}..."
	rsync -avz --delete --exclude='node_modules' --exclude='.git' . /opt/app/

deploy-to-remote:
	@echo "🚀 Deploying ${APP_NAME}..."
	rsync -avz --exclude='node_modules' --exclude='.git' . amijar-vm:/opt/app/

test-api:
	@echo "📬 Sample curl commands for testing the API"
	@echo ""
	@echo "👉 Create a new review:"
	@echo "curl -X POST http://localhost:${PORT}/reviews \\"
	@echo "     -H 'Content-Type: application/json' \\"
	@echo "     -d '{\"reviewText\": \"Great product!\", \"approved\": false}'"
	@echo ""
	@echo "👉 List all reviews:"
	@echo "curl http://localhost:${PORT}/reviews"
	@echo ""
	@echo "👉 Get a review by ID (replace <id>):"
	@echo "curl http://localhost:${PORT}/reviews/<id>"
	@echo ""
	@echo "👉 Approve a review (replace <id>):"
	@echo "curl -X PUT http://localhost:${PORT}/reviews/<id> \\"
	@echo "     -H 'Content-Type: application/json' \\"
	@echo "     -d '{\"approved\": true}'"
	@echo ""
	@echo "👉 Update review text and approval (replace <id>):"
	@echo "curl -X PUT http://localhost:${PORT}/reviews/<id> \\"
	@echo "     -H 'Content-Type: application/json' \\"
	@echo "     -d '{\"reviewText\": \"Updated text\", \"approved\": true}'"
	@echo ""
	@echo "👉 Delete a review by ID (replace <id>):"
	@echo "curl -X DELETE http://localhost:${PORT}/reviews/<id>"
	@echo ""
	@echo "🧪 Tip: Use jq to pretty-print: | jq"


#------------------------------------
# use for later development ONLY
dev:
	@echo "🚀 Starting server with nodemon..."
	nodemon $(ENTRY)

dev-mongo-start:
	@echo "🔋 Starting MongoDB using Homebrew (macOS only)..."
	brew services start mongodb-community

dev-mongo-stop:
	@echo "🛑 Stopping MongoDB using Homebrew (macOS only)..."
	brew services stop mongodb-community

dev-mongo-status:
	brew services list | grep mongodb


#------------------------------------
# FOR APP SERVICE - USE ONLY ON LINUX VM
mongo-install:
	sudo apt update -y
	sudo apt install -y mongodb
	sudo systemctl start mongodb
	sudo systemctl enable mongodb

mongo-uninstall:
	sudo systemctl stop mongodb
	sudo systemctl disable mongodb
	sudo apt remove -y mongodb

mongo-status:
	sudo systemctl status mongodb

service-install:
	sudo cp systemd/app.service /etc/systemd/system/app.service
	sudo systemctl daemon-reload
	sudo systemctl enable app
	sudo systemctl start app
	sudo systemctl enable app

service-uninstall:
	sudo systemctl stop app
	sudo systemctl disable app
	sudo rm /etc/systemd/system/app.service
	sudo systemctl daemon-reload

service-down:
	sudo systemctl stop app

service-up:
	sudo systemctl start app

service-restart:
	sudo systemctl restart app

service-status:
	sudo systemctl status app
	# sudo journalctl -u app -f

logger-setup:
	sudo touch /var/log/app.log
	sudo chown ${USER}:www-data /var/log/app.log
	sudo chmod a+w /var/log/app.log


#------------------------------------
# SITE-SPECIFIC CONVENIENCE TARGETS
# These wrap the generic targets from pub-site-setup.mk
#------------------------------------

# OilYourHair.com specific targets
.PHONY: oilyourhair-deploy oilyourhair-conf-deploy oilyourhair-html-deploy
.PHONY: oilyourhair-status oilyourhair-logs

oilyourhair-deploy:
	@$(MAKE) site-deploy SITE=$(OILYOURHAIR_SITE)

oilyourhair-conf-deploy:
	@$(MAKE) site-conf-deploy SITE=$(OILYOURHAIR_SITE)

oilyourhair-html-deploy:
	@$(MAKE) site-html-deploy SITE=$(OILYOURHAIR_SITE)

oilyourhair-status:
	@$(MAKE) site-status SITE=$(OILYOURHAIR_SITE)

oilyourhair-logs:
	@$(MAKE) site-logs SITE=$(OILYOURHAIR_SITE)

oilyourhair-tunnel-create:
	@$(MAKE) site-tunnel-create SITE=$(OILYOURHAIR_SITE)

oilyourhair-tunnel-config:
	@$(MAKE) site-tunnel-config SITE=$(OILYOURHAIR_SITE)

oilyourhair-tunnel-route:
	@$(MAKE) site-tunnel-route SITE=$(OILYOURHAIR_SITE)

oilyourhair-tunnel-test:
	@$(MAKE) site-tunnel-test SITE=$(OILYOURHAIR_SITE)

#------------------------------------
# CLOUDFLARE TUNNEL GLOBAL SETUP
#------------------------------------
.PHONY: cloudflare-install cloudflare-login cloudflare-dns-credentials cloudflare-setup

cloudflare-install:
	@echo "📦 Installing cloudflared..."
	wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
	sudo dpkg -i cloudflared-linux-amd64.deb
	rm cloudflared-linux-amd64.deb
	@echo "✅ cloudflared installed"

cloudflare-login:
	@echo "🔐 Logging into Cloudflare..."
	cloudflared tunnel login
	@echo "✅ Login complete"

cloudflare-dns-credentials:
	@echo "🔑 Setting up Cloudflare DNS credentials for certbot..."
	@mkdir -p ~/.secrets
	@read -p "Enter your Cloudflare API Token: " token; \
	echo "dns_cloudflare_api_token = $$token" > ~/.secrets/cloudflare.ini
	@chmod 600 ~/.secrets/cloudflare.ini
	@echo "✅ Credentials saved to ~/.secrets/cloudflare.ini"

cloudflare-setup:
	@echo "🚀 Complete Cloudflare Tunnel Setup for $(OILYOURHAIR_SITE)"
	@echo ""
	@echo "Step 1: Install cloudflared"
	@$(MAKE) cloudflare-install
	@echo ""
	@echo "Step 2: Login to Cloudflare"
	@$(MAKE) cloudflare-login
	@echo ""
	@echo "Step 3: Create tunnel"
	@$(MAKE) oilyourhair-tunnel-create
	@echo ""
	@echo "Step 4: Configure tunnel"
	@$(MAKE) oilyourhair-tunnel-config
	@echo ""
	@echo "✅ Setup complete! Next steps:"
	@echo "   1. Test tunnel: make oilyourhair-tunnel-test"
	@echo "   2. Route DNS: make oilyourhair-tunnel-route"
	@echo "   3. Install as service: make site-tunnel-service SITE=$(OILYOURHAIR_SITE)"

cloudflare-tunnel-service:
	@echo "🔧 Installing Cloudflare tunnel as systemd service..."
	sudo cloudflared service install
	sudo systemctl start cloudflared
	sudo systemctl enable cloudflared
	@echo "✅ Tunnel service installed and started"

cloudflare-tunnel-status:
	sudo systemctl status cloudflared

cloudflare-tunnel-logs:
	sudo journalctl -u cloudflared -f

#------------------------------------
# LEGACY TARGETS (kept for backwards compatibility)
# Consider migrating to site-specific targets above
#------------------------------------

frontend-conf-deploy: oilyourhair-conf-deploy
	@echo "⚠️  'frontend-conf-deploy' is deprecated, use 'oilyourhair-conf-deploy'"

frontend-deploy: oilyourhair-html-deploy
	@echo "⚠️  'frontend-deploy' is deprecated, use 'oilyourhair-html-deploy'"

