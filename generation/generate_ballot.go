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

// generateBallot creates a single ballot that has range proofs for each selection along with a
// range proof for the adherence to vote limits.
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

// generateSelectionsForContest creates all selections + range proofs for each selection for a given contest along with
// a range proof for adhering to the vote limits.
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

		// Generating range proof based on El Gamal encryption, the vote, and the selection limit
		rangeProof := generateRangeProofFromEncryptionAndNonce(*alpha, *beta, *epsilon, publicKey, selectionLimit, m)

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
