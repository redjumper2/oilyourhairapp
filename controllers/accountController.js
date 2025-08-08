// accountController.js

const {
  getAllAccounts,
  getAccountById,
  addAccount,
  updateAccount,
  deleteAccount,
} = require('../models/accountModel');

async function fetchAllAccounts(req, res) {
  const accounts = await getAllAccounts();
  res.json(accounts);
}

async function fetchAccountById(req, res) {
  try {
    const account = await getAccountById(req.params.id);
    if (!account) return res.status(404).json({ error: 'Account not found' });
    res.json(account);
  } catch {
    res.status(400).json({ error: 'Invalid ID format' });
  }
}

async function createAccount(req, res) {
  const { email, name } = req.body;

  if (typeof email !== 'string') {
    return res.status(400).json({ error: 'email is required and must be a string' });
  }
  if (typeof name !== 'string') {
    return res.status(400).json({ error: 'name is required and must be a string' });
  }
  const account = await addAccount(`${email}`, `${name}`);
  res.status(201).json(account);
}

async function editAccount(req, res) {
  try {
    const { email, name } = req.body;
    const updateObj = {};
    if (email !== undefined) updateObj.email = email;
    if (name !== undefined) updateObj.name = name;

    const updated = await updateAccount(req.params.id, updateObj);
    if (!updated) return res.status(404).json({ error: 'Account not found' });
    res.json(updated);
  } catch {
    res.status(400).json({ error: 'Invalid ID format' });
  }
}

async function removeAccount(req, res) {
  try {
    const success = await deleteAccount(req.params.id);
    if (!success) return res.status(404).json({ error: 'Account not found' });
    res.json({ message: 'Deleted successfully' });
  } catch {
    res.status(400).json({ error: 'Invalid ID format' });
  }
}

module.exports = {
  fetchAllAccounts,
  fetchAccountById,
  createAccount,
  editAccount,
  removeAccount,
};
