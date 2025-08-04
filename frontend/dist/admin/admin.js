// app.js - This is the main entry point for the OilYourHair.com App.
// http://yourdomain.com/ → Frontend
// http://yourdomain.com/api/reviews → API

const API_URL = '/api/reviews';
const reviewsContainer = document.getElementById('reviews');
const form = document.getElementById('review-form');

async function loadReviews() {
  const res = await fetch(API_URL);
  const data = await res.json();

  reviewsContainer.innerHTML = data.map(r => `
    <div class="review">
      <p><strong>ID:</strong> ${r._id}</p>
      <p><strong>Rating:</strong> ${r.rating}</p>
      <p><strong>Text:</strong> ${r.reviewText}</p>
      <p><strong>Approved:</strong> ${r.approved}</p>
      <button onclick="approveReview('${r._id}')">✅ Approve</button>
      <button onclick="deleteReview('${r._id}')">🗑 Delete</button>
    </div>
  `).join('');
}

form.addEventListener('submit', async (e) => {
  e.preventDefault();
  // Default to 3 If .value is falsy (meaning it’s an empty string "", null, undefined, or 0) 
  const rating = document.getElementById('rating').value || 5;
  const reviewText = document.getElementById('reviewText').value;
  const approved = document.getElementById('approved').checked;

  await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ rating, reviewText, approved })
  });

  form.reset();
  loadReviews();
});

async function approveReview(id) {
  await fetch(`${API_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ approved: true })
  });
  loadReviews();
}

async function deleteReview(id) {
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
  loadReviews();
}

// Initial load
loadReviews();