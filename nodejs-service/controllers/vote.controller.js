import { prisma } from "../configs/prisma.js";

export const castVote = async (req, res) => {
  try {
    const { voter_id, candidate_id } = req.body;
    const result = await prisma.$transaction(async (tx) => {
      const newVote = await tx.vote.create({
        data: {
          voter_id,
          candidate_id,
        },
        select: {
          vote_id: true,
          voter_id: true,
          candidate_id: true,
          createdAt: true
        }
      });
      const candidate = await prisma.candidate.findFirst({ where: { candidate_id } });
      const voter = await prisma.voter.findFirst({ where: { voter_id } });
      let weight;
      if (voter.age >= 18 && voter.age < 30) weight = 1;
      else if (voter.age >= 30 && voter.age < 40) weight = 2;
      else if (voter.age >= 40 && voter.age < 50) weight = 3;
      else if (voter.age > 50) weight = 4;
      await prisma.candidate.update({
        where: {
          candidate_id,
        },
        data: {
          votes_count: candidate.votes_count + 1,
        }
      })
      return { newVote };
    })

    res.status(200).json(result.newVote);
  } catch (error) {
    console.log(error);
  }
}

export const getCandidateVotes = async (req, res) => {
  try {
    const { candidate_id } = req.params;
    const candidateVotes = await prisma.vote.findMany({
      where: {
        candidate_id: parseInt(candidate_id),
      },
      select: {
        candidate_id: true,
      }
    });
    const candidate = await prisma.candidate.findFirst({ where: { candidate_id: parseInt(candidate_id) } })
    res.status(200).json({ candidate_id, votes: candidate.votes_count });
  } catch (error) {
    console.log(error);
  }
}

export const filterCandidatesByParty = async (req, res) => {
  try {
    const { party_name } = req.params;
    const filteredCandidates = await prisma.candidate.findMany({
      where: {
        party: party_name
      }
    });
    res.status(200).json(filteredCandidates);
  } catch (error) {
    console.log(error);
  }
}

export const votingResults = async (req, res) => {
  try {
    const results = await prisma.candidate.findMany({
      orderBy: {
        votes_count: "desc"
      }
    });
    res.status(200).json({ results });
  } catch (error) {
    console.log(error);
  }
}


export const winningCandidate = async (req, res) => {
  try {
    const maxVotesCandidate = await prisma.candidate.aggregate({
      _max: {
        votes_count: true
      }
    });
    const maxVotes = maxVotesCandidate._max.votes_count;
    const winners = await prisma.candidate.findMany({
      where: {
        votes_count: maxVotes
      },
      select: {
        candidate_id: true,
        name: true,
        votes_count: true
      },
      orderBy: {
        name: 'asc'
      }
    });
    res.status(200).json({ winners });
  } catch (error) {
    console.log(error);
  }
}