package generation

import (
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	mod "electionguard-sandbox-go/modular_arithmetic"
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
		ballotContests[i] = generateBallotContest(contest, publicKey)
	}

	ballot := models.Ballot{Contests: ballotContests}
	return ballot
}

// generateBallotContest creates all selections + range proofs for each selection for a given contest along with
// a range proof for adhering to the vote limits.
func generateBallotContest(contest models.Contest, publicKey crypto.PublicKey) models.BallotContest {
	selectionLimit := contest.SelectionLimit
	contestSelectionLimit := contest.ContestSelectionLimit

	// Creating list of selections for ballot/contest combination
	amountOfSelections := len(contest.Selections)
	selections := make([]models.BallotSelection, amountOfSelections)

	// Nonces and encryption values that are valid (i.e. within the selection limit)
	encryptionNonces := getRandomNumbersModQ(amountOfSelections)
	encryptionValues := getRandomVotes(amountOfSelections, selectionLimit, contestSelectionLimit)

	alphaHat := big.NewInt(1)
	betaHat := big.NewInt(1)
	epsilonHat := big.NewInt(0)

	for k, selection := range contest.Selections {
		// Get (message, nonce) and generate El Gamal encryption
		m := encryptionValues[k]
		epsilon := encryptionNonces[k]
		alpha, beta := crypto.Encrypt(publicKey, m, epsilon)

		// Calculating the product of all encryptions / sum of nonces
		alphaHat = mod.MulP(alphaHat, alpha)
		betaHat = mod.MulP(betaHat, beta)
		epsilonHat = mod.AddQ(epsilonHat, epsilon)

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

	// Generate vote limit range proof (adherence to vote limit)
	encryptionValuesSum := 0
	for _, value := range encryptionValues {
		encryptionValuesSum += value
	}
	voteAdherenceRangeProof := generateRangeProofFromEncryptionAndNonce(*alphaHat, *betaHat, *epsilonHat, publicKey, contestSelectionLimit, encryptionValuesSum) // note that we always encrypt "selection limit" thus no undervotes is performed

	// Creating ballot contest
	ballotContest := models.BallotContest{
		BallotSelections: selections,
		Proof:            voteAdherenceRangeProof,
	}

	return ballotContest
}

func getRandomVotes(amountOfSelections, selectionLimit, contestSelectionLimit int) []int {
	votes := make([]int, amountOfSelections)
	sumOfVotes := 0

	for i := 0; i < amountOfSelections; i++ {
		vote := rand.Intn(selectionLimit + 1)
		for vote+sumOfVotes > contestSelectionLimit {
			vote = rand.Intn(selectionLimit + 1)
		}

		votes[i] = vote
		sumOfVotes += vote
	}

	return votes
}

func getRandomNumbersModQ(n int) []*big.Int {
	nonces := make([]*big.Int, n)

	for i := 0; i < n; i++ {
		nonces[i] = crypto.GenerateRandomModQ()
	}

	return nonces
}
