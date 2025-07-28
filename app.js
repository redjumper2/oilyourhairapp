// app.js
require('dotenv').config();

const port = process.env.PORT || 3000;
const express = require('express');
const { connectToMongo } = require('./db/mongo');
const reviewRoutes = require('./routes/reviewRoutes');
const logger = require('./utils/logger');

const app = express();

app.use(express.json());
app.use('/reviews', reviewRoutes);

async function start() {
  try {
    await connectToMongo();
    app.listen(port, () => {
      logger.info(`🚀 Server running at http://localhost:${port}`);
      logger.info(`testing logger output`);
    });
  } catch (err) {
    logger.error(`❌ Server startup error: ${err.message}`);
  }
}

start();
