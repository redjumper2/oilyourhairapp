// review.js

const { ObjectId } = require('mongodb');

class Review {
  constructor({ _id = null, rating, reviewText, approved = false }) {
    if (!reviewText) throw new Error('reviewText is required');
    this._id = _id ? new ObjectId(_id) : undefined;
    this.rating = rating;
    this.reviewText = reviewText;
    this.approved = approved;
  }

  static fromMongo(doc) {
    return new Review({
      _id: doc._id,
      rating: doc.rating,
      reviewText: doc.reviewText,
      approved: doc.approved
    });
  }

  toMongo() {
    return {
      rating: this.rating,
      reviewText: this.reviewText,
      approved: this.approved
    };
  }
}

module.exports = Review;
