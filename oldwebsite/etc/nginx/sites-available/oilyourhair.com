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

    location /admin {
        auth_basic "Restricted Access";
        auth_basic_user_file /etc/nginx/.htpasswd;

    }


    # Proxy API calls
    # location /api/ {
    #     proxy_pass http://localhost:3000/;
    #     proxy_http_version 1.1;
    #     proxy_set_header Upgrade $http_upgrade;
    #     proxy_set_header Connection 'upgrade';
    #     proxy_set_header Host $host;
    #     proxy_cache_bypass $http_upgrade;
    # }

    # Optional: deny access to hidden files (like .git)
    location ~ /\. {
        deny all;
    }


    listen 8080;
    listen [::]:8080;
}
