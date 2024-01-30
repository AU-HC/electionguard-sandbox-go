package generation

import (
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	"math/big"
	"math/rand"
)

func GenerateBallots(manifest models.Manifest, amountOfBallots int, publicKey crypto.PublicKey) []models.Ballot {
	ballots := make([]models.Ballot, amountOfBallots)

	for i := 0; i < amountOfBallots; i++ {
		ballot := generateBallot(manifest, publicKey)
		ballots[i] = ballot
	}

	return ballots
}

func generateBallot(manifest models.Manifest, publicKey crypto.PublicKey) models.Ballot {
	ballotContests := make([]models.BallotContest, len(manifest.Contests))

	for i, contest := range manifest.Contests {
		ballotSelectionsForContest := generateSelectionsForContest(contest, publicKey)
		ballotContest := models.BallotContest{
			BallotSelections:       ballotSelectionsForContest,
			Proof:                  models.RangeProof{},
			CiphertextAccumulation: models.Ciphertext{},
		}
		ballotContests[i] = ballotContest
	}

	ballot := models.Ballot{
		Contests: ballotContests,
	}
	return ballot
}

func generateSelectionsForContest(contest models.Contest, publicKey crypto.PublicKey) []models.BallotSelection {
	// Creating ballot contest
	selectionLimit := contest.SelectionLimit

	// Creating list of selections for ballot/contest combination
	amountOfSelections := len(contest.Selections)
	selections := make([]models.BallotSelection, amountOfSelections)

	encryptionNonces := getEncryptionNonces(amountOfSelections)
	encryptionValues := getRandomVotes(amountOfSelections, selectionLimit)

	for k, selection := range contest.Selections {
		// Get (message, nonce) and generate El Gamal encryption
		m := encryptionValues[k]
		epsilon := encryptionNonces[k]
		alpha, beta := crypto.Encrypt(publicKey, m, epsilon)

		// Generate random challenges
		cpProofs := make([]models.ChaumPedersenProof, selectionLimit+1) // need R+1 proofs
		challenges := getEncryptionNonces(selectionLimit + 1)           // need R+1 challenges
		commitments := getEncryptionNonces(selectionLimit + 1)          // need R+1 commitments

		// Generate range proof of this encryption
		for j := 0; j <= selectionLimit; j++ {
			u := commitments[j]
			c := *challenges[j]

			var a, b big.Int
			a.Exp(publicKey.G, u, publicKey.P)
			if m == j {
				b.Exp(publicKey.K, u, publicKey.P)
			} else {
				var cMul big.Int
				t := new(big.Int)

				t = t.Add(u, cMul.Mul(big.NewInt(int64(selectionLimit-j)), &c))
				t.Mod(t, publicKey.Q)
				b.Exp(publicKey.K, t, publicKey.P)
			}

			// Calculating the response
			var v big.Int
			v.Sub(u, epsilon.Mul(epsilon, &c))
			v.Mod(&v, publicKey.Q)

			// Filling the values in we have calculated (note that the c for m == j has to be replaced later)
			cpProofs[j] = models.ChaumPedersenProof{
				Challenge:     c,
				ProofPad:      a,
				ProofData:     b,
				ProofResponse: v,
			}
		}

		// Calculating "true" claim proof
		c := crypto.HMAC(*publicKey.K, 0x21, publicKey.K, alpha, beta)
		cl := new(big.Int)
		cl = cl.Set(c)
		for j := 0; j <= selectionLimit; j++ {
			if m != j {
				cl.Sub(cl, challenges[j])
				cl.Mod(cl, publicKey.Q)
			}
		}
		var v big.Int
		v.Sub(commitments[m], epsilon.Mul(epsilon, c))
		v.Mod(&v, publicKey.Q)
		cpProofs[m] = models.ChaumPedersenProof{
			Challenge:     *c,
			ProofPad:      cpProofs[m].ProofPad,
			ProofData:     cpProofs[m].ProofData,
			ProofResponse: v,
		}

		// Saving all the proofs into range proof struct
		rangeProof := models.RangeProof{
			Challenge:  *c,
			Proofs:     cpProofs,
			RangeLimit: selectionLimit,
		}

		// Creating ballot selection
		ballotSelection := models.BallotSelection{
			ObjectId:   selection.Name,
			Ciphertext: models.Ciphertext{Pad: *alpha, Data: *beta},
			Proof:      rangeProof,
		}

		// Setting the ballot selection to the correct index
		selections[k] = ballotSelection
	}

	// Generate vote limit range proof
	// ...

	return selections
}

func getRandomVotes(amountOfSelections int, limit int) []int {
	votes := make([]int, amountOfSelections)
	sumOfVotes := 0

	for i := 0; i < amountOfSelections; i++ {
		vote := rand.Intn(limit - sumOfVotes)
		votes[i] = vote
		sumOfVotes += vote
	}

	return votes
}

func getEncryptionNonces(n int) []*big.Int {
	nonces := make([]*big.Int, n)

	for i := 0; i < n; i++ {
		nonces[i] = crypto.GenerateRandomModQ()
	}

	return nonces
}
