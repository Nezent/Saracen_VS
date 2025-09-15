import express from 'express';
import cors from "cors";
import dotenv from "dotenv";

dotenv.config();
const app = express();
app.use(cors());
app.use(express.json());
const PORT = process.env.PORT || 8000;


app.get("/", (req, res) => {
  res.send("Hackathon API is runnings");
});



app.listen(PORT, "0.0.0.0", () => {
  console.log(`Server is dancing on http://localhost:${PORT} \n${new Date(Date.now()).toLocaleTimeString()}
    `);
});