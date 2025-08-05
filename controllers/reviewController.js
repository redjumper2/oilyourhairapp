// reviewController.js

const {
  getAllReviews,
  getReviewById,
  addReview,
  updateReview,
  deleteReview,
} = require('../models/reviewModel');

async function fetchAll(req, res) {
  const reviews = await getAllReviews();
  res.json(reviews);
}

async function fetchById(req, res) {
  try {
    const review = await getReviewById(req.params.id);
    if (!review) return res.status(404).json({ error: 'Review not found' });
    res.json(review);
  } catch {
    res.status(400).json({ error: 'Invalid ID format' });
  }
}

async function create(req, res) {
  const { user, rating, reviewText, approved } = req.body;

  if (typeof reviewText !== 'string') {
    return res.status(400).json({ error: 'reviewText is required and must be a string' });
  }
  const review = await addReview(`${user}`, +rating, reviewText, !!approved);
  res.status(201).json(review);
}

async function update(req, res) {
  try {
    const { reviewText, approved } = req.body;
    const updateObj = {};
    if (reviewText !== undefined) updateObj.reviewText = reviewText;
    if (approved !== undefined) updateObj.approved = approved;

    const updated = await updateReview(req.params.id, updateObj);
    if (!updated) return res.status(404).json({ error: 'Review not found' });
    res.json(updated);
  } catch {
    res.status(400).json({ error: 'Invalid ID format' });
  }
}

async function remove(req, res) {
  try {
    const success = await deleteReview(req.params.id);
    if (!success) return res.status(404).json({ error: 'Review not found' });
    res.json({ message: 'Deleted successfully' });
  } catch {
    res.status(400).json({ error: 'Invalid ID format' });
  }
}

module.exports = {
  fetchAll,
  fetchById,
  create,
  update,
  remove,
};
