// review.js

const { ObjectId } = require('mongodb');

class Account {
  constructor({ _id = null, email, name}) {
    if (!reviewText) throw new Error('reviewText is required');
    this._id = _id ? new ObjectId(_id) : undefined;
    this.email = email;
    this.name = name;
  }

  static fromMongo(doc) {
    return new Review({
      _id: doc._id,
      email: doc.email,
      name: doc.name
    });
  }

  toMongo() {
    return {
      email: this.email,
      name: this.name
    };
  }
}

module.exports = Review;
