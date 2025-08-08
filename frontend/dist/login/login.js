// login.js — script for adding accounts
const API_URL = '/api/accounts'; // change to your real accounts API endpoint
const accountsContainer = document.getElementById('accounts');
const form = document.getElementById('account-form');

async function loadAccounts() {
  const res = await fetch(API_URL);
  const data = await res.json();

  accountsContainer.innerHTML = data.map(acc => `
    <div class="account">
      <p><strong>ID:</strong> ${acc._id}</p>
      <p><strong>Email:</strong> ${acc.email}</p>
      <p><strong>Name:</strong> ${acc.name}</p>
      <button onclick="deleteAccount('${acc._id}')">🗑 Delete</button>
    </div>
  `).join('');
}

form.addEventListener('submit', async (e) => {
  e.preventDefault();

  const email = document.getElementById('email').value;
  const name = document.getElementById('name').value;

  await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, name })
  });

  form.reset();
  loadAccounts();
});

async function deleteAccount(id) {
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
  loadAccounts();
}

// Initial load
loadAccounts();
