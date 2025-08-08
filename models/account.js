// review.js

const { ObjectId } = require('mongodb');

class Account {
  constructor({ _id = null, email, name}) {
    // validating the input parameters
    if (!email || email.trim() === '') {
      throw new Error('email cannot be empty');
    }
    if (!name || name.trim() === '') {
      throw new Error('name cannot be empty');
    }

    this._id = _id ? new ObjectId(_id) : undefined;
    this.email = email;
    this.name = name;
  }

  static fromMongo(doc) {
    return new Account({
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

module.exports = Account;
