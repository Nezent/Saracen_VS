import { prisma } from "../configs/prisma.js";

export const registerCandidate = async (req, res) => {
  const { name, party } = req.body;
  try {
    if (!name || !party) {
      res.status(401).json({
        success: false,
        message: "missing fields",
      });
    }
    const newCandidate = await prisma.candidate.create({
      data: {
        name,
        party
      }
    });
    res.status(200).json(newCandidate);
  } catch (error) {
    console.log(error);
  }
}

export const listCandidates = async (req, res) => {
  try {
    const {party} = req.query;
    let candidates;
    if(party) {
      candidates = await prisma.candidate.findMany({
        where: {
          party: party
        }
      })
    } else candidates = await prisma.candidate.findMany();
    return res.status(200).json(candidates);
  } catch (error) {
    console.log(error);
  }
}