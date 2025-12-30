FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Create directory structure
RUN mkdir -p /var/www/oilyourhair.com/html \
    && mkdir -p /etc/nginx/sites-available \
    && mkdir -p /etc/nginx/sites-enabled

# Copy nginx configuration files
COPY sites-available/oilyourhair.com /etc/nginx/sites-available/oilyourhair.com
COPY sites-enabled/oilyourhair.com /etc/nginx/sites-enabled/oilyourhair.com

# Update main nginx.conf to include sites-enabled
RUN sed -i '/http {/a \    include /etc/nginx/sites-enabled/*;' /etc/nginx/nginx.conf

# Note: Website files will be mounted via volume in docker-compose

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
