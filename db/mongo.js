const { MongoClient } = require('mongodb');

const uri = process.env.MONGO_URI || 'mongodb://localhost:27017';
const dbName = process.env.MONGO_DB || 'oilyourhairdb';

let db;

async function connectToMongo() {
  const client = new MongoClient(uri, { useUnifiedTopology: true });
  await client.connect();
  console.log('✅ Connected to MongoDB');
  db = client.db(dbName);

  const collection = db.collection('reviews');
  const count = await collection.countDocuments();
  if (count === 0) {
    await collection.insertOne({ user: 'test user', rating: 5, reviewText: 'First test review', approved: false });
  }

  const collection2 = db.collection('accounts');
  const count2 = await collection.countDocuments();
  if (count2 === 0) {
    await collection2.insertOne({ email: 'example@gmail.com', name: 'Test User' });
  }
}

function getDb() {
  if (!db) throw new Error('❌ DB not initialized.');
  return db;
}

module.exports = { connectToMongo, getDb };
