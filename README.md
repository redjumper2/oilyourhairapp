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
# list
curl -XGET http://localhost:3000/reviews

curl -XGET https://oilyourhair.com/api/reviews
[{"_id":"688766bc764eebb390c89b6e","reviewText":"test","approved":false}]

# get
curl -XGET http://localhost:3000/reviews/688766bc764eebb390c89b6e

curl -XGET https://oilyourhair.com/api/reviews/688766bc764eebb390c89b6e

response:
{"_id":"688766bc764eebb390c89b6e","reviewText":"testing2","approved":false}

# add - 
curl -XPOST -H "Content-Type: application/json"  http://localhost:3000/reviews -d '{"user":"test user", "rating": 5, "reviewText": "Love it!", "approved": false}'

curl -XPOST -H "Content-Type: application/json"  https://oilyourhair.com/api/reviews -d '{"user":"test user", "rating": 5, "reviewText": "Love it!", "approved": false}'

# list - 
curl -XGET https://localhost:3000/reviews
[{"_id":"6887675c764eebb390c89b70","reviewText":"testing2","approved":false},{"_id":"688767a1764eebb390c89b71","reviewText":"Love it!","approved":false}]

curl -XGET https://oilyourhair.com/api/reviews
[{"_id":"6887675c764eebb390c89b70","reviewText":"testing2","approved":false},{"_id":"688767a1764eebb390c89b71","reviewText":"Love it!","approved":false}]

# delete review - curl -XDELETE https://oilyourhair.com/api/reviews/6887675c764eebb390c89b70
curl -XDELETE https://localhost:3000/reviews/6887675c764eebb390c89b70

curl -XDELETE https://oilyourhair.com/api/reviews/6887675c764eebb390c89b70

# Update
curl -X PUT http://localhost:3000/reviews/688767a1764eebb390c89b71 -H "Content-Type: application/json" -d '{"approved": true}'

curl -X PUT https://oilyourhair.com/api/reviews/688767a1764eebb390c89b71 -H "Content-Type: application/json" -d '{"approved": true}'
```
