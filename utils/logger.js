require('dotenv').config();
const { createLogger, format, transports } = require('winston');

const logger = createLogger({
  level: process.env.LOG_LEVEL || 'info',  // <--- uses .env level
  format: format.combine(
    format.timestamp(),
    format.printf(({ timestamp, level, message }) => {
      return `${timestamp} [${level.toUpperCase()}] ${message}`;
    })
  ),
  transports: [
    new transports.File({ filename: '/var/log/app.log' })
  ],
});

// Optional: Add console logging in development
if (process.env.NODE_ENV !== 'production') {
  logger.add(new transports.Console());
}

module.exports = logger;

// if daily rotation is needed, install winston-daily-rotate-file
// npm install winston-daily-rotate-file

// const DailyRotateFile = require('winston-daily-rotate-file');

// transports: [
//   new DailyRotateFile({
//     filename: '/var/log/app-%DATE%.log',
//     datePattern: 'YYYY-MM-DD',
//     maxFiles: '14d'
//   }),
//   new winston.transports.Console()
// ]
