#!/bin/bash

# Script to seed sample reviews into the database

API_URL="http://localhost:8002/api/v1/public/oilyourhair.com/reviews"

echo "Seeding sample reviews..."

# Review 1 - Sarah Johnson
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Nourishing Hair Oil",
    "name": "Sarah Johnson",
    "rating": 5,
    "text": "Absolutely love this argan oil! My hair has never been softer or shinier. I use it after every wash and the results are incredible. The light coconut scent is amazing too!",
    "highlight": "Best hair product I have ever used!"
  }'

echo ""

# Review 2 - Michael Chen
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Hydrating Hair Serum",
    "name": "Michael Chen",
    "rating": 5,
    "text": "This serum is a game-changer for my dry, damaged hair. After just two weeks of use, my hair feels healthier and looks more vibrant. Highly recommend!",
    "highlight": "Transformed my damaged hair completely"
  }'

echo ""

# Review 3 - Emily Rodriguez
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Moisturizing Hair Mask",
    "name": "Emily Rodriguez",
    "rating": 5,
    "text": "I have curly hair that tends to get frizzy, but this hair mask keeps it smooth and defined. The natural ingredients are a huge plus, and it smells wonderful!",
    "highlight": "Perfect for curly hair - no more frizz!"
  }'

echo ""

# Review 4 - David Thompson
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Strengthening Hair Oil",
    "name": "David Thompson",
    "rating": 5,
    "text": "As someone with thinning hair, I was skeptical, but this oil has made a noticeable difference. My hair feels thicker and stronger. Will definitely keep using this!",
    "highlight": "Noticeable improvement in hair thickness"
  }'

echo ""

# Review 5 - Jessica Park
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Repair Hair Treatment",
    "name": "Jessica Park",
    "rating": 5,
    "text": "After years of heat styling and coloring, my hair was in bad shape. This treatment has brought it back to life! It is soft, shiny, and so much healthier. Cannot recommend enough!",
    "highlight": "Brought my damaged hair back to life"
  }'

echo ""

# Review 6 - Ryan Martinez
curl -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{
    "product": "Shine & Smooth Hair Oil",
    "name": "Ryan Martinez",
    "rating": 5,
    "text": "I have been using this oil for a month now and the shine it gives my hair is incredible. It is lightweight and does not make my hair greasy. A little goes a long way!",
    "highlight": "Amazing shine without the grease"
  }'

echo ""
echo "✅ All sample reviews seeded successfully!"
