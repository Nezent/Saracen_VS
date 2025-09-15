import express from 'express';
import { listCandidates, registerCandidate } from '../controllers/candidate.controller.js';
import { getCandidateVotes } from '../controllers/vote.controller.js';

const router = express.Router();

router.post('/', registerCandidate);
router.get('/', listCandidates);
router.get('/:candidate_id/votes', getCandidateVotes);


export default router;