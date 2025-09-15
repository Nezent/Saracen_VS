# Saracen Voting System 🗳️

A comprehensive voting system implementation featuring advanced ballot types including encrypted ballots, ranked-choice voting with Schulze method, and weighted voting mechanisms.

## 🚀 Quick Start

### Prerequisites
- Docker and Docker Compose installed on your system
- Git (to clone the repository)

### Installation & Setup

1. **Download the project**
   ```bash
   git clone <repository-url>
   cd Saracen
   ```

2. **Run with Docker Compose**
   ```bash
   docker-compose up
   ```

3. **Access the API**
   The server will be running on **port 8000**:
   ```
   http://localhost:8000
   ```

4. **Health Check**
   Verify the server is running:
   ```bash
   curl http://localhost:8000/health
   ```

## 📋 Available APIs

### Voter Management (Q1-Q5)
- `POST /api/voters` - Create a new voter
- `GET /api/voters/{voter_id}` - Get voter information  
- `GET /api/voters` - List all voters
- `PUT /api/voters/{voter_id}` - Update voter information
- `DELETE /api/voters/{voter_id}` - Delete a voter

### Vote Operations (Q13-Q15)
- `GET /api/votes/timeline?candidate_id={id}` - Get vote timeline for candidate
- `POST /api/votes/weighted` - Cast a weighted vote
- `GET /api/votes/range?candidate_id={id}&from={t1}&to={t2}` - Get votes in time range

### Advanced Ballot Systems
- `POST /api/ballots/encrypted` - Submit encrypted ballot (Q16)
- `POST /api/ballots/ranked` - Submit ranked-choice ballot (Q19)
- `GET /api/ballots/ranked/results?election_id={id}` - Get Schulze method results

## 🧪 Quick API Tests

### Create a Voter
```bash
curl -X POST http://localhost:8000/api/voters \
-H "Content-Type: application/json" \
-d '{"voter_id": 1, "name": "Alice", "age": 25}'
```

### Cast a Weighted Vote  
```bash
curl -X POST http://localhost:8000/api/votes/weighted \
-H "Content-Type: application/json" \
-d '{"voter_id": 1, "candidate_id": 2}'
```

### Submit Encrypted Ballot
```bash
curl -X POST http://localhost:8000/api/ballots/encrypted \
-H "Content-Type: application/json" \
-d '{
  "election_id": "nat-2025",
  "voter_id": 100,
  "ciphertext": "my_cipher_text",
  "zk_proof": "my_proof",
  "voter_pubkey": "1",
  "nullifier": "unique123",
  "signature": "my_signature"
}'
```

### Submit Ranked Ballot
```bash
curl -X POST http://localhost:8000/api/ballots/ranked \
-H "Content-Type: application/json" \
-d '{
  "election_id": "city-rcv-2025",
  "voter_id": 123,
  "ranking": [3, 1, 2],
  "timestamp": "2025-09-15T10:15:00Z"
}'
```

## 🏗️ Architecture

This system implements **Clean Architecture** with the following layers:

- **Domain Layer**: Business models and rules
- **Application Layer**: Use cases and business logic  
- **Infrastructure Layer**: Database and external services
- **Interface Layer**: HTTP handlers and API endpoints

## 🔧 Features

- ✅ **Basic Voter Management**: CRUD operations with validation
- ✅ **Weighted Voting**: Vote weight based on voter profile activity
- ✅ **Time-based Queries**: Vote timeline and range queries
- ✅ **Encrypted Ballots**: Zero-knowledge proof support with nullifier validation
- ✅ **Ranked Choice Voting**: Schulze method implementation for winner determination
- ✅ **Database Integration**: PostgreSQL with proper foreign key constraints
- ✅ **Input Validation**: Automatic base64/hex conversion for cryptographic fields

## 🗄️ Database

The system uses PostgreSQL with the following main tables:
- `voter` - Voter information and voting status
- `candidate` - Candidate details and vote counts  
- `votes` - Individual votes with weights and timestamps
- `encrypted_ballots` - Encrypted ballot submissions with proofs
- `ranked_ballots` & `ballot_rankings` - Ranked choice voting data

## 📖 API Documentation

For complete API documentation with request/response examples, see the `api_list.json` file which contains detailed specifications for all 20 problem sets (Q1-Q20).

## 🛠️ Development

Built with:
- **Go 1.23** - Backend service
- **PostgreSQL** - Database
- **Docker** - Containerization
- **Gorilla Mux** - HTTP routing

## 📝 Notes

- The system auto-converts simple inputs to proper formats (base64 for cryptographic data, hex for keys)
- Weighted voting uses profile update activity to determine vote weights
- Ranked choice voting implements the Schulze method for determining winners
- All cryptographic validations are simplified for development purposes

---