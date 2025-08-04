# Project settings
APP_NAME = user-reviews-app
ENTRY = app.js
PORT = 3000
MONGO_URI = mongodb://localhost:27017
MONGO_DB = myappdb

# Default target
.PHONY: help
help:
	@echo "📦 ${APP_NAME} Makefile"
	@echo ""
	@echo "Commands:"
	@echo "  install        Install npm dependencies"
	@echo "  run            Run the server"
	@echo "  dev            Run with nodemon (if installed)"
	@echo "  mongo-start    Start MongoDB (macOS: brew only)"
	@echo "  mongo-stop     Stop MongoDB (macOS: brew only)"
	@echo "  mongo-status   Show MongoDB status"
	@echo "  test-api       Show test curl commands"


#------------------------------------
# use for development
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
# USE ONLY ON LINUX VM
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
# use for production deployment
# nginx and certbot installation and configuration
nginx-install:
	sudo apt update -y
	sudo apt install nginx -y

certbot-install:
	sudo apt install certbot python3-certbot-nginx -y

certbot-renew:
	sudo certbot renew --nginx
	sudo certbot --nginx -d oilyourhair.com -d oilyourhair.com
	sudo systemctl status certbot.timer

frontend-deploy:
	sudo cp nginx-app-conf/oilyourhair.com /etc/nginx/sites-available/oilyourhair.com
	sudo ln -sf /etc/nginx/sites-available/oilyourhair.com /etc/nginx/sites-enabled/
	sudo nginx -t
	sudo systemctl reload nginx
