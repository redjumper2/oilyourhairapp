// login.js - Handles login-related actions for OilYourHair.com
// http://yourdomain.com/ → Frontend
// http://yourdomain.com/api/account → API

const API_URL = '/api/account';
const accountsContainer = document.getElementById('accounts');
const form = document.getElementById('account-form');

async function loadAccounts() {
  const res = await fetch(API_URL);
  const data = await res.json();

  accountsContainer.innerHTML = '';
  data.forEach(renderAccount); // use the reusable function
}

function renderAccount(account) {
  const div = document.createElement('div');
  div.classList.add('account');
  div.innerHTML = `
    <p><strong>ID:</strong> ${account._id}</p>
    <p><strong>Email:</strong> ${account.email}</p>
    <p><strong>Name:</strong> ${account.name}</p>
    <button onclick="deleteAccount('${account._id}', this)">🗑 Delete</button>
  `;
  accountsContainer.appendChild(div);
}

form.addEventListener('submit', async (e) => {
  e.preventDefault();
  const email = document.getElementById('email').value;
  const name = document.getElementById('name').value || 'Anonymous';

  await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, name })
  });

  form.reset();
  loadAccounts(); // Refresh list
});

async function deleteAccount(id, button) {
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
  loadAccounts(); // Refresh after deletion
}

// Initial load
loadAccounts();
