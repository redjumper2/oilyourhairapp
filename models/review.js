const { ObjectId } = require('mongodb');

class Review {
  constructor({ _id = null, reviewText, approved = false }) {
    if (!reviewText) throw new Error('reviewText is required');
    this._id = _id ? new ObjectId(_id) : undefined;
    this.reviewText = reviewText;
    this.approved = approved;
  }

  static fromMongo(doc) {
    return new Review({
      _id: doc._id,
      reviewText: doc.reviewText,
      approved: doc.approved
    });
  }

  toMongo() {
    return {
      reviewText: this.reviewText,
      approved: this.approved
    };
  }
}

module.exports = Review;
