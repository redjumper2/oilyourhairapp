# Description
An app for OilYourHair

# Development stack
MEN stack (MongoDB, Express, Node, Plain/Simple/Original html/js)

# to deploy from dev to production
```bash
make deploy
```

# add to git
```bash
git add <filenames> or '.'
git commit -m "<comment>"
Optional but ideal - [git push]
```

# starting and stopping service and viewing logs
```bash
make service-down
make service-up

make service-status

# to view its logs
tail -f /var/log/app.log
```

# Testing with curl
```bash
# Create
curl -X POST http://localhost:3000/reviews -H "Content-Type: application/json" -d '{"reviewText": "Love it!", "approved": false}'

# List all
curl http://localhost:3000/reviews

# Get by ID
curl http://localhost:3000/reviews/<id>

# Update
curl -X PUT http://localhost:3000/reviews/<id> -H "Content-Type: application/json" -d '{"approved": true}'

# Delete
curl -X DELETE http://localhost:3000/reviews/<id>
```
