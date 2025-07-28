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
    await collection.insertOne({ reviewText: 'First test review', approved: false });
  }
}

function getDb() {
  if (!db) throw new Error('❌ DB not initialized.');
  return db;
}

module.exports = { connectToMongo, getDb };
