/*
  Warnings:

  - You are about to drop the `Candidate` table. If the table is not empty, all the data it contains will be lost.
  - You are about to drop the `users` table. If the table is not empty, all the data it contains will be lost.

*/
-- DropForeignKey
ALTER TABLE "public"."votes" DROP CONSTRAINT "votes_candidate_id_fkey";

-- DropForeignKey
ALTER TABLE "public"."votes" DROP CONSTRAINT "votes_voter_id_fkey";

-- DropTable
DROP TABLE "public"."Candidate";

-- DropTable
DROP TABLE "public"."users";

-- CreateTable
CREATE TABLE "public"."voter" (
    "voter_id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "age" INTEGER NOT NULL,
    "has_voted" BOOLEAN NOT NULL,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "voter_pkey" PRIMARY KEY ("voter_id")
);

-- CreateTable
CREATE TABLE "public"."candidate" (
    "candidate_id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "party" TEXT NOT NULL,
    "votes_count" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "candidate_pkey" PRIMARY KEY ("candidate_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "voter_voter_id_key" ON "public"."voter"("voter_id");

-- CreateIndex
CREATE UNIQUE INDEX "candidate_candidate_id_key" ON "public"."candidate"("candidate_id");

-- AddForeignKey
ALTER TABLE "public"."votes" ADD CONSTRAINT "votes_voter_id_fkey" FOREIGN KEY ("voter_id") REFERENCES "public"."voter"("voter_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "public"."votes" ADD CONSTRAINT "votes_candidate_id_fkey" FOREIGN KEY ("candidate_id") REFERENCES "public"."candidate"("candidate_id") ON DELETE RESTRICT ON UPDATE CASCADE;
