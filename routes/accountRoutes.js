
// accountRoutes.js
const express = require('express');
const {
  fetchAllAccounts,
  fetchAccountById,
  createAccount,
  editAccount,
  removeAccount,
} = require('../controllers/accountController');

const router = express.Router();

router.get('/', fetchAllAccounts);
router.get('/:id', fetchAccountById);
router.post('/', createAccount);
router.put('/:id', editAccount);
router.delete('/:id', removeAccount);

module.exports = router;
