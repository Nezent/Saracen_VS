-- Migration: Add encrypted_ballots and ranked_ballots tables for Q16 and Q19
-- Created: 2025-09-15 14:00:00

-- CreateTable for Q16: Encrypted Ballots
CREATE TABLE "public"."encrypted_ballots" (
    "ballot_id" TEXT NOT NULL,
    "election_id" TEXT NOT NULL,
    "voter_id" INTEGER NOT NULL,
    "ciphertext" TEXT NOT NULL,
    "zk_proof" TEXT NOT NULL,
    "voter_pubkey" TEXT NOT NULL,
    "nullifier" TEXT NOT NULL,
    "signature" TEXT NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'accepted',
    "anchored_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "encrypted_ballots_pkey" PRIMARY KEY ("ballot_id")
);

-- CreateTable for Q19: Ranked Ballots
CREATE TABLE "public"."ranked_ballots" (
    "ballot_id" TEXT NOT NULL,
    "election_id" TEXT NOT NULL,
    "voter_id" INTEGER NOT NULL,
    "timestamp" TIMESTAMP(3) NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'accepted',

    CONSTRAINT "ranked_ballots_pkey" PRIMARY KEY ("ballot_id")
);

-- CreateTable for Q19: Ballot Rankings (individual candidate rankings)
CREATE TABLE "public"."ballot_rankings" (
    "id" SERIAL NOT NULL,
    "ballot_id" TEXT NOT NULL,
    "candidate_id" INTEGER NOT NULL,
    "rank_position" INTEGER NOT NULL,

    CONSTRAINT "ballot_rankings_pkey" PRIMARY KEY ("id")
);

-- CreateIndex for Q16: Unique nullifier to prevent double voting
CREATE UNIQUE INDEX "encrypted_ballots_nullifier_key" ON "public"."encrypted_ballots"("nullifier");

-- CreateIndex for Q19: Optimize queries by ballot_id
CREATE INDEX "ballot_rankings_ballot_id_idx" ON "public"."ballot_rankings"("ballot_id");

-- CreateIndex for Q19: Optimize queries by candidate_id
CREATE INDEX "ballot_rankings_candidate_id_idx" ON "public"."ballot_rankings"("candidate_id");

-- AddForeignKey for Q16: Link encrypted ballots to voters
ALTER TABLE "public"."encrypted_ballots" ADD CONSTRAINT "encrypted_ballots_voter_id_fkey" 
FOREIGN KEY ("voter_id") REFERENCES "public"."voter"("voter_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey for Q19: Link ranked ballots to voters
ALTER TABLE "public"."ranked_ballots" ADD CONSTRAINT "ranked_ballots_voter_id_fkey" 
FOREIGN KEY ("voter_id") REFERENCES "public"."voter"("voter_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey for Q19: Link ballot rankings to ranked ballots
ALTER TABLE "public"."ballot_rankings" ADD CONSTRAINT "ballot_rankings_ballot_id_fkey" 
FOREIGN KEY ("ballot_id") REFERENCES "public"."ranked_ballots"("ballot_id") ON DELETE CASCADE;

-- AddForeignKey for Q19: Link ballot rankings to candidates
ALTER TABLE "public"."ballot_rankings" ADD CONSTRAINT "ballot_rankings_candidate_id_fkey" 
FOREIGN KEY ("candidate_id") REFERENCES "public"."candidate"("candidate_id") ON DELETE RESTRICT ON UPDATE CASCADE;