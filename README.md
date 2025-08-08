# Description
An app for OilYourHair

# Objective
[o] Create a webstie that allows the brand, "oilyourhair.com" to showcase, sell, and market their product.

# Product Roadmap


Phase 1 (MVP):

Goal: Get the brand pages out, get testimonials/reviews

[x] main website pages html/css

[x] admin/reviews crud backend

[x] admin page for review approval

[x] password protected for review approval

[ ] implement crud for Account(s) (user)

[ ] implement approved reviews page

[ ] Take website live

---
Phase 2 (Manage Sales):

[ ] product list backend with pricing

[ ] shopping cart backend

[ ] shopping cart frontend

[ ] Inventory management

---
Phase 3 (Manage Marketing - internal):

[ ] email creation templates for campaign

[ ] email distribution for campaign

---
Phase 4 (Role management for full website mgmt):

[ ] implement crud for Role(s)

[ ] implement role control (3 roles - end user, admin, review moderator)

[ ] implement google authN for Account authN/authZ with role(s)

[ ] admin page - add more functionality - review approval, email creation using brand template, email distribution, promotions, discounts


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

curl -XPOST -H "Content-Type: application/json"  https://oilyourhair.com/api/accounts -d '{"email":"example@gmail.com", "name": "Test User"}'

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
