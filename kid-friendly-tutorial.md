# 👦 Build a Review API (Kid-Friendly Tutorial)

## 🧠 What Are We Making?

We’re making a cool program where people can leave reviews — like for books, games, or snacks! This is called a **Review API**.

## 🧰 Tools We’ll Use

- **Node.js** – our robot helper that runs code
- **Express.js** – lets us talk to our app with URLs
- **MongoDB** – where we store reviews
- **curl/Postman** – to test our review system

## 🛠️ Step-by-Step Guide

### 1. 🎯 What’s an API?

An **API** is like a restaurant waiter. You tell it what you want (GET or POST), and it gives you the food (data). Simple!

- `GET` – Get reviews
- `POST` – Add a review
- `PUT` – Update a review
- `DELETE` – Delete a review

### 2. 📦 Project Structure

```
my-review-app/
├── app.js                 # Main app file
├── routes/
│   └── reviewRoutes.js    # URL routes
├── controllers/
│   └── reviewController.js # Logic for routes
├── models/
│   ├── review.js          # Review data shape
│   └── reviewModel.js     # Database code
├── db/
│   └── mongo.js           # MongoDB connection
├── utils/
│   └── logger.js          # Logging to file
```

### 3. 🚀 Start the App

```bash
npm install
node app.js
```

You should see:

```
🚀 Server running at http://localhost:3000
```

### 4. 🧪 Try with curl

Get all reviews:

```bash
curl http://localhost:3000/reviews
```

Add a review:

```bash
curl -X POST http://localhost:3000/reviews \
  -H "Content-Type: application/json" \
  -d '{"reviewText": "This is awesome!"}'
```

### 5. 🌐 Or Use Postman (Easier)

- Download Postman
- Make a new POST request to:

```
http://localhost:3000/reviews
```

- Set Body > raw > JSON:

```json
{
  "reviewText": "Amazing!"
}
```

Click Send. Yay! 🎉

### 6. 🔁 Update a Review

```bash
curl -X PUT http://localhost:3000/reviews/<id> \
  -H "Content-Type: application/json" \
  -d '{"approved": true}'
```

### 7. ❌ Delete a Review

```bash
curl -X DELETE http://localhost:3000/reviews/<id>
```

## 🕵️‍♂️ Where It Logs

Check the file:

```
/var/log/app.log
```

It shows messages when you add/update/delete reviews.

---

## 🎥 Watch These First

### ✅ What is an API? (5 min)

[https://www.youtube.com/watch?v=Q-BpqyOT3a8](https://www.youtube.com/watch?v=Q-BpqyOT3a8)

### ✅ curl Explained (5 min)

[https://www.youtube.com/watch?v=Fs5bGktWzPY](https://www.youtube.com/watch?v=Fs5bGktWzPY)

### ✅ Postman Beginner Tutorial

[https://www.youtube.com/watch?v=VywxIQ2ZXw4](https://www.youtube.com/watch?v=VywxIQ2ZXw4)

---

## ✨ Wrap-Up

You just built your own API! You're a mini software engineer now. 🚀 Want to try building one for pizza reviews next? 🍕

