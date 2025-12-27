# Generic Public Site Setup Makefile
# Include this in your main Makefile
# Usage: make target SITE=example.com

# Required variables (pass via command line or set in main Makefile):
# SITE - domain name (e.g., oilyourhair.com)

#------------------------------------
# Directory paths
NGINX_CONF_SRC = ./etc/nginx/sites-available/$(SITE)
NGINX_CONF_DEST = /etc/nginx/sites-available/$(SITE)
NGINX_ENABLED = /etc/nginx/sites-enabled/$(SITE)
WWW_SRC = ./var/www/$(SITE)/html
WWW_DEST = /var/www/$(SITE)/html

#------------------------------------
# Nginx setup
.PHONY: site-nginx-install

site-nginx-install:
	sudo apt update -y
	sudo apt install nginx -y

#------------------------------------
# Site-specific deployment
.PHONY: site-conf-deploy site-html-deploy site-deploy

site-conf-deploy:
	@echo "🔧 Deploying nginx config for $(SITE)..."
	sudo cp $(NGINX_CONF_SRC) $(NGINX_CONF_DEST)
	sudo ln -sf $(NGINX_CONF_DEST) $(NGINX_ENABLED)
	sudo nginx -t
	sudo systemctl reload nginx
	@echo "✅ Nginx config deployed for $(SITE)"

site-html-deploy:
	@echo "📦 Deploying HTML for $(SITE)..."
	sudo mkdir -p $(WWW_DEST)
	sudo rsync -av --exclude='admin' $(WWW_SRC)/ $(WWW_DEST)/
	sudo chown -R www-data:www-data $(WWW_DEST)
	@echo "✅ HTML deployed for $(SITE)"

site-deploy: site-conf-deploy site-html-deploy
	@echo "✅ Full deployment complete for $(SITE)"

#------------------------------------
# Cloudflare Tunnel setup (per-site)
.PHONY: site-tunnel-create site-tunnel-config site-tunnel-route

site-tunnel-create:
	@echo "🚇 Creating Cloudflare Tunnel for $(SITE)..."
	cloudflared tunnel create $(SITE)-tunnel
	@echo "✅ Tunnel created for $(SITE)"

site-tunnel-config:
	@echo "⚙️  Creating tunnel config for $(SITE)..."
	@mkdir -p ~/.cloudflared
	@TUNNEL_ID=$$(cloudflared tunnel list | grep "$(SITE)-tunnel" | awk '{print $$1}'); \
	CONFIG_FILE=~/.cloudflared/$(SITE)-config.yml; \
	echo "tunnel: $$TUNNEL_ID" > $$CONFIG_FILE; \
	echo "credentials-file: /home/$(USER)/.cloudflared/$$TUNNEL_ID.json" >> $$CONFIG_FILE; \
	echo "" >> $$CONFIG_FILE; \
	echo "ingress:" >> $$CONFIG_FILE; \
	echo "  - hostname: $(SITE)" >> $$CONFIG_FILE; \
	echo "    service: http://localhost:80" >> $$CONFIG_FILE; \
	echo "  - hostname: www.$(SITE)" >> $$CONFIG_FILE; \
	echo "    service: http://localhost:80" >> $$CONFIG_FILE; \
	echo "  - service: http_status:404" >> $$CONFIG_FILE
	@echo "✅ Config created at $$CONFIG_FILE"

site-tunnel-route:
	@echo "🌐 Routing DNS for $(SITE) to tunnel..."
	cloudflared tunnel route dns $(SITE)-tunnel $(SITE)
	cloudflared tunnel route dns $(SITE)-tunnel www.$(SITE)
	@echo "✅ DNS routed for $(SITE)"

site-tunnel-test:
	@echo "🧪 Testing tunnel for $(SITE) (Ctrl+C to stop)..."
	cloudflared tunnel --config ~/.cloudflared/$(SITE)-config.yml run $(SITE)-tunnel

#------------------------------------
# Site removal/cleanup
.PHONY: site-remove

site-remove:
	@echo "🗑️  Removing $(SITE) configuration..."
	@read -p "Are you sure you want to remove $(SITE)? [y/N] " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		sudo rm -f $(NGINX_ENABLED); \
		sudo rm -f $(NGINX_CONF_DEST); \
		sudo nginx -t && sudo systemctl reload nginx; \
		echo "✅ $(SITE) nginx config removed"; \
	else \
		echo "❌ Cancelled"; \
	fi

#------------------------------------
# Utilities
.PHONY: site-status site-logs

site-status:
	@echo "📊 Status for $(SITE):"
	@echo "  Nginx config: $(NGINX_CONF_DEST)"
	@ls -lh $(NGINX_CONF_DEST) 2>/dev/null || echo "    ❌ Not found"
	@echo "  Nginx enabled: $(NGINX_ENABLED)"
	@ls -lh $(NGINX_ENABLED) 2>/dev/null || echo "    ❌ Not found"
	@echo "  Web root: $(WWW_DEST)"
	@ls -lh $(WWW_DEST) 2>/dev/null || echo "    ❌ Not found"
	@echo "  Cloudflare Tunnel:"
	@cloudflared tunnel info $(SITE)-tunnel 2>/dev/null || echo "    ❌ Not found"

site-logs:
	@echo "📜 Nginx logs for $(SITE):"
	@if [ -f /var/log/nginx/$(SITE).access.log ]; then \
		echo "Access log:"; \
		sudo tail -20 /var/log/nginx/$(SITE).access.log; \
	fi
	@if [ -f /var/log/nginx/$(SITE).error.log ]; then \
		echo ""; \
		echo "Error log:"; \
		sudo tail -20 /var/log/nginx/$(SITE).error.log; \
	fi
