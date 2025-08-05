server {

    server_name oilyourhair.com www.oilyourhair.com;

    root /var/www/oilyourhair.com/html;
    index index.html index.htm;

    # Logging
    access_log /var/log/nginx/oilyourhair.access.log;
    error_log  /var/log/nginx/oilyourhair.error.log;

    # Serve static files
    location / {
        add_header Cache-Control "no-cache, no-store, must-revalidate";
        add_header Pragma "no-cache";
        add_header Expires 0;

        try_files $uri $uri/ =404;
    }

#    location ~* \.js$ {
#        # Prevent aggressive caching
#        add_header Cache-Control "no-cache, no-store, must-revalidate";
#        add_header Pragma "no-cache";
#        add_header Expires 0;
#    }

#    location ~* \.html$ {
#        # Prevent aggressive caching
#        add_header Cache-Control "no-cache, no-store, must-revalidate";
#        add_header Pragma "no-cache";
#        add_header Expires 0;
#    }

    # Proxy API calls
    location /api/ {
        proxy_pass http://localhost:3000/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Optional: deny access to hidden files (like .git)
    location ~ /\. {
        deny all;
    }

    listen [::]:443 ssl ipv6only=on; # managed by Certbot
    listen 443 ssl; # managed by Certbot
    ssl_certificate /etc/letsencrypt/live/oilyourhair.com/fullchain.pem; # managed by Certbot
    ssl_certificate_key /etc/letsencrypt/live/oilyourhair.com/privkey.pem; # managed by Certbot
    include /etc/letsencrypt/options-ssl-nginx.conf; # managed by Certbot
    ssl_dhparam /etc/letsencrypt/ssl-dhparams.pem; # managed by Certbot


}

server {
    if ($host = www.oilyourhair.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    if ($host = oilyourhair.com) {
        return 301 https://$host$request_uri;
    } # managed by Certbot


    listen 80;
    listen [::]:80;

    server_name oilyourhair.com www.oilyourhair.com;
    return 404; # managed by Certbot




}