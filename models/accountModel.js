// accountModel.js

const { ObjectId } = require('mongodb');
const { getDb } = require('../db/mongo');
const Account = require('./account');
const logger = require('../utils/logger');

function getAccountCollection() {
  return getDb().collection('accounts');
}

async function getAllAccounts() {
  const docs = await getAccountCollection().find().toArray();
  logger.info(`Fetched ${docs.length} accounts`);
  logger.info('Accounts viewed successfully!');
  return docs.map(Account.fromMongo);
}

async function getAccountById(id) {
  const account = await getAccountCollection().findOne({ _id: new ObjectId(id) });
  logger.info(`Fetched account by ID: ${id}`);
  return account ? Account.fromMongo(account) : null;
}

async function addAccount(email, name = 'Anonymous') {
  const account = new Account({ email, name });
  const result = await getAccountCollection().insertOne(account.toMongo());
  logger.info(`Added account with ID: ${result.insertedId}`);
  return new Account({ _id: result.insertedId, ...account });
}

async function updateAccount(id, update) {
  const result = await getAccountCollection().findOneAndUpdate(
    { _id: new ObjectId(id) },
    { $set: update },
    { returnDocument: 'after' }
  );
  logger.info(`Updated account with ID: ${id}`);
  return result.value;
}

async function deleteAccount(id) {
  const result = await getAccountCollection().deleteOne({ _id: new ObjectId(id) });
  logger.info(`Deleted account with ID: ${id}`);
  return result.deletedCount === 1;
}

module.exports = {
  getAllAccounts,
  getAccountById,
  addAccount,
  updateAccount,
  deleteAccount,
};
