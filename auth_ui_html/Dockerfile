FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy static files
COPY *.html *.js *.css /usr/share/nginx/html/

# Fix file permissions
RUN chmod 644 /usr/share/nginx/html/*.html /usr/share/nginx/html/*.js /usr/share/nginx/html/*.css

# Create simple nginx config
RUN echo 'server { \
    listen 5173; \
    server_name _; \
    root /usr/share/nginx/html; \
    index login.html; \
    \
    location / { \
        try_files $uri $uri.html $uri/ =404; \
    } \
    \
    # Rewrite /login to /login.html (preserve query string) \
    location = /login { \
        rewrite ^ /login.html break; \
    } \
    \
    # Rewrite /invite to /invite.html (preserve query string) \
    location = /invite { \
        rewrite ^ /invite.html break; \
    } \
    \
    # CORS headers for API calls \
    add_header Access-Control-Allow-Origin * always; \
}' > /etc/nginx/conf.d/default.conf

EXPOSE 5173

CMD ["nginx", "-g", "daemon off;"]
