FROM nginx:alpine

# Remove default nginx site
RUN rm -f /etc/nginx/conf.d/default.conf

# Expose HTTP port
EXPOSE 80

# Note: nginx.conf and site configs are mounted via volumes in docker-compose.yml

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
