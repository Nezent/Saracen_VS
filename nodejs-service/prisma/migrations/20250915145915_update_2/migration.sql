-- CreateTable
CREATE TABLE "public"."ranked_ballots" (
    "ballot_id" TEXT NOT NULL,
    "election_id" TEXT NOT NULL,
    "voter_id" INTEGER NOT NULL,
    "timestamp" TIMESTAMP(3) NOT NULL,
    "status" TEXT NOT NULL DEFAULT 'accepted',

    CONSTRAINT "ranked_ballots_pkey" PRIMARY KEY ("ballot_id")
);

-- CreateTable
CREATE TABLE "public"."ballot_rankings" (
    "id" SERIAL NOT NULL,
    "ballot_id" TEXT NOT NULL,
    "candidate_id" INTEGER NOT NULL,
    "rank_position" INTEGER NOT NULL,

    CONSTRAINT "ballot_rankings_pkey" PRIMARY KEY ("id")
);

-- CreateTable
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

-- CreateIndex
CREATE UNIQUE INDEX "encrypted_ballots_nullifier_key" ON "public"."encrypted_ballots"("nullifier");

-- AddForeignKey
ALTER TABLE "public"."ranked_ballots" ADD CONSTRAINT "ranked_ballots_voter_id_fkey" FOREIGN KEY ("voter_id") REFERENCES "public"."voter"("voter_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "public"."ballot_rankings" ADD CONSTRAINT "ballot_rankings_ballot_id_fkey" FOREIGN KEY ("ballot_id") REFERENCES "public"."ranked_ballots"("ballot_id") ON DELETE CASCADE ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "public"."ballot_rankings" ADD CONSTRAINT "ballot_rankings_candidate_id_fkey" FOREIGN KEY ("candidate_id") REFERENCES "public"."candidate"("candidate_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "public"."encrypted_ballots" ADD CONSTRAINT "encrypted_ballots_voter_id_fkey" FOREIGN KEY ("voter_id") REFERENCES "public"."voter"("voter_id") ON DELETE RESTRICT ON UPDATE CASCADE;
