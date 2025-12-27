// reviewModel.js

const { ObjectId } = require('mongodb');
const { getDb } = require('../db/mongo');
const Review = require('./review');
const logger = require('../utils/logger');

function getReviewCollection() {
  return getDb().collection('reviews');
}

async function getAllReviews() {
  const docs = await getReviewCollection().find().toArray();
  logger.info(`Fetched ${docs.length} reviews`);
  logger.info('Reviews viewed successfully!');
  return docs.map(Review.fromMongo);
}

async function getReviewById(id) {
  const review = await getReviewCollection().findOne({ _id: new ObjectId(id) });
  logger.info(`Fetched review by ID: ${id}`);
  return review ? Review.fromMongo(review) : null;
}

async function addReview(user = 'Anonymous', rating = 5, reviewText, approved = false) {
  const review = new Review({ user, rating, reviewText, approved });
  const result = await getReviewCollection().insertOne(review.toMongo());
  logger.info(`Added review with ID: ${result.insertedId}`);
  return new Review({ _id: result.insertedId, ...review });
}

async function updateReview(id, update) {
  const result = await getReviewCollection().findOneAndUpdate(
    { _id: new ObjectId(id) },
    { $set: update },
    { returnDocument: 'after' }
  );
  logger.info(`Updated review with ID: ${id}`);
  return result.value;
}

async function deleteReview(id) {
  const result = await getReviewCollection().deleteOne({ _id: new ObjectId(id) });
  logger.info(`Deleted review with ID: ${id}`);
  return result.deletedCount === 1;
}

module.exports = {
  getAllReviews,
  getReviewById,
  addReview,
  updateReview,
  deleteReview,
};
