package verification

import (
	"electionguard-sandbox-go/models"
	"fmt"
)

func verifyStep5(ballots []models.BallotSelection) {
	for _, ballot := range ballots {
		fmt.Println(ballot)
	}
}
