import express from 'express';
import { votingResults, winningCandidate } from '../controllers/vote.controller.js';

const router = express.Router();

router.get('/', votingResults);
router.get('/winner', winningCandidate);

export default router;