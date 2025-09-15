/*
  Warnings:

  - The primary key for the `users` table will be changed. If it partially fails, the table could be left without primary key constraint.
  - You are about to drop the column `email` on the `users` table. All the data in the column will be lost.
  - You are about to drop the column `id` on the `users` table. All the data in the column will be lost.
  - You are about to drop the column `password` on the `users` table. All the data in the column will be lost.
  - A unique constraint covering the columns `[voter_id]` on the table `users` will be added. If there are existing duplicate values, this will fail.
  - Added the required column `age` to the `users` table without a default value. This is not possible if the table is not empty.
  - Added the required column `has_voted` to the `users` table without a default value. This is not possible if the table is not empty.

*/
-- DropIndex
DROP INDEX "public"."users_email_key";

-- AlterTable
ALTER TABLE "public"."users" DROP CONSTRAINT "users_pkey",
DROP COLUMN "email",
DROP COLUMN "id",
DROP COLUMN "password",
ADD COLUMN     "age" INTEGER NOT NULL,
ADD COLUMN     "has_voted" BOOLEAN NOT NULL,
ADD COLUMN     "voter_id" SERIAL NOT NULL,
ADD CONSTRAINT "users_pkey" PRIMARY KEY ("voter_id");

-- CreateTable
CREATE TABLE "public"."Candidate" (
    "candidate_id" SERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "party" TEXT NOT NULL,
    "votes_count" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "Candidate_pkey" PRIMARY KEY ("candidate_id")
);

-- CreateTable
CREATE TABLE "public"."votes" (
    "vote_id" SERIAL NOT NULL,
    "voter_id" INTEGER NOT NULL,
    "candidate_id" INTEGER NOT NULL,
    "weight" INTEGER NOT NULL DEFAULT 1,
    "created_at" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "votes_pkey" PRIMARY KEY ("vote_id")
);

-- CreateIndex
CREATE UNIQUE INDEX "Candidate_candidate_id_key" ON "public"."Candidate"("candidate_id");

-- CreateIndex
CREATE UNIQUE INDEX "votes_vote_id_key" ON "public"."votes"("vote_id");

-- CreateIndex
CREATE UNIQUE INDEX "users_voter_id_key" ON "public"."users"("voter_id");

-- AddForeignKey
ALTER TABLE "public"."votes" ADD CONSTRAINT "votes_voter_id_fkey" FOREIGN KEY ("voter_id") REFERENCES "public"."users"("voter_id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "public"."votes" ADD CONSTRAINT "votes_candidate_id_fkey" FOREIGN KEY ("candidate_id") REFERENCES "public"."Candidate"("candidate_id") ON DELETE RESTRICT ON UPDATE CASCADE;
