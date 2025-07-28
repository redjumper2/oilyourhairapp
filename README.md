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
curl -XGET https://oilyourhair.com/api/reviews
[{"_id":"688766bc764eebb390c89b6e","reviewText":"test","approved":false}]

# get
curl -XGET https://oilyourhair.com/api/reviews/688766bc764eebb390c89b6e
{"_id":"688766bc764eebb390c89b6e","reviewText":"testing2","approved":false}

# add
curl -XPOST -H "Content-Type: application/json"  https://oilyourhair.com/api/reviews -d '{"reviewText": "Love it!", "approved": false}'

# list
curl -XGET https://oilyourhair.com/api/reviews
[{"_id":"6887675c764eebb390c89b70","reviewText":"testing2","approved":false},{"_id":"688767a1764eebb390c89b71","reviewText":"Love it!","approved":false}]

# delete review
curl -XDELETE https://oilyourhair.com/api/reviews/6887675c764eebb390c89b70

# Update
curl -X PUT http://localhost:3000/api/reviews/688767a1764eebb390c89b71 -H "Content-Type: application/json" -d '{"approved": true}'
```
