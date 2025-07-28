const express = require('express');
const {
  fetchAll,
  fetchById,
  create,
  update,
  remove,
} = require('../controllers/reviewController');

const router = express.Router();

router.get('/', fetchAll);
router.get('/:id', fetchById);
router.post('/', create);
router.put('/:id', update);
router.delete('/:id', remove);

module.exports = router;
