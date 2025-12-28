# Description
This is a simplified specification for Reviews API.

# API specification
## Create Review
POST /reviews
Body:
{
  "reviewText":"some review text",
  "approved": false
}

Response: 201 OK (response created)

# Get reviews
GET /api/reviews

Response: 200 OK
Body:
[{"_id":"123", "reviewText":"some review", "approved": false}, {"_id":"123", "reviewText":"some review", "approved": false}]

# Get 1 review
GET /api/review

Response: 200 OK
Body:
{{"_id":"123", "reviewText":"some review", "approved": false}}

# Delete 1 review
DELETE /api/review/123

Response: 200 OK
{"message":"Deleted successfully"}

Response: 404
{"error":"Review not found"}

# Put 1 review
PUT /api/review/123

Response: 200 OK
