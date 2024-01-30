package generation

import (
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	"math/big"
)

func GenerateBallots(manifest models.Manifest, amountOfBallots int, publicKey crypto.PublicKey) []models.Ballot {
	ballots := make([]models.Ballot, amountOfBallots)

	for i := 0; i < amountOfBallots; i++ {
		ballots[i] = generateBallot(manifest)
	}

	return ballots
}

func generateBallot(manifest models.Manifest) models.Ballot {
	ballot := models.Ballot{}
	ballotContests := make([]models.BallotContest, len(manifest.Contests))

	for i, contest := range manifest.Contests {
		// Creating ballot contest
		ballotContest := models.BallotContest{ObjectId: contest.Name}

		// Creating list of selections for ballot/contest combination
		amountOfSelections := len(contest.Selections)
		selections := make([]models.BallotSelection, amountOfSelections)
		encryptionNonces := getEncryptionNonces(amountOfSelections)
		selectionLimit := contest.SelectionLimit

		for _, selections := range contest.Selections {
			// Generate encryption of 'x' < selectionLimit?

			// Generate range proof of this encryption

		}
		// Generate vote limit range proof
		// ...

		// Adding ballot selections to ballot contest and appending the contest to the contests list
		ballotContest.BallotSelections = selections
		ballotContests[i] = ballotContest
	}

	return ballot
}

func getEncryptionNonces(n int) []*big.Int {
	nonces := make([]*big.Int, n)

	for i := 0; i < n; i++ {
		nonces[i] = crypto.GenerateRandomModQ()
	}

	return nonces
}
