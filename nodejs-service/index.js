import express from 'express';
import cors from "cors";
import dotenv from "dotenv";
import axios from 'axios';
import candidateRoutes from './routes/candidate.routes.js';
import voteRoutes from "./routes/vote.routes.js";
import resultRoutes from "./routes/results.routes.js";

const GOLANG_API_URL = process.env.GOLANG_API_URL;

dotenv.config();
const app = express();
app.use(cors());
app.use(express.json());
const PORT = process.env.PORT || 8000;


/// -------------- golang apis ---------------
app.post('/api/voters', async (req, res) => {
  const data = req.body;
  const response = await axios.post(`${GOLANG_API_URL}/api/voters`, data);
  res.status(200).json(response.data);
});
app.get('/api/voters', async (req, res) => {
  const response = await axios.get(`${GOLANG_API_URL}/api/voters`);
  res.status(200).json(response.data);
});
app.get('/api/voters/:voter_id', async (req, res) => {
  const { voter_id } = req.params;
  const response = await axios.get(`${GOLANG_API_URL}/api/voters/${voter_id}`);
  res.status(200).json(response.data);
});
app.put('/api/voters/:voter_id', async (req, res) => {
  const { voter_id } = req.params;
  const data = req.body;
  const response = await axios.put(`${GOLANG_API_URL}/api/voters/${voter_id}`, data);
  res.status(200).json(response.data);
});
app.delete('/api/voters/:voter_id', async (req, res) => {
  const { voter_id } = req.params;
  const response = await axios.delete(`${GOLANG_API_URL}/api/voters/${voter_id}`);
  res.status(200).json(response.data);
});


/// ----------- 13 - 20 -------------
app.get('/api/votes/timeline', async (req, res) => { //13
  const { candidate_id } = req.params;
  const response = await axios.get(`${GOLANG_API_URL}/api/votes/timeline?candidate_id=${candidate_id}`);
  res.status(200).json(response.data);
});
app.post('/api/votes/weighted', async (req, res) => { //14
  const response = await axios.post(`${GOLANG_API_URL}/api/votes/weighted`);
  res.status(200).json(response.data);
});
app.get('/api/votes/range', async (req, res) => { //15
  const { candidate_id, t1, t2 } = req.query;
  const response = await axios.get(`${GOLANG_API_URL}/api/votes/range?candidate_id=${candidate_id}}&from=${t1}&to=${t2}`);
  res.status(200).json(response.data);
});
app.post('/api/ballots/encrypted', async (req, res) => { //16
  const data = req.body;
  const response = await axios.post(`${GOLANG_API_URL}/api/ballots/encrypted`, data);
  res.status(200).json(response.data);
});


app.post('/api/ballots/ranked', async (req, res) => { //19
  const data = req.body;
  const response = await axios.post(`${GOLANG_API_URL}/api/ballots/ranked`, data);
  res.status(200).json(response.data);
})


app.use('/api/candidates', candidateRoutes);
app.use('/api/votes', voteRoutes);
app.use('/api/results', resultRoutes);
app.get("/", (req, res) => {
  res.send("Hackathon API is runnings");
});



app.listen(PORT, "0.0.0.0", () => {
  console.log(`Server is dancing on http://localhost:${PORT} \n${new Date(Date.now()).toLocaleTimeString()}
    `);
});