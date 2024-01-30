package verification

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	"fmt"
	"math/big"
)

func VerifyStep6(ballots []models.Ballot, publicKey crypto.PublicKey) {
	g := publicKey.G
	k := publicKey.K

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			alphaHat, betaHat := computeContestTotal(contest.BallotSelections)

			toBeHashed := []interface{}{*publicKey.K, alphaHat, betaHat}
			computedC := big.NewInt(0)
			for j, proof := range contest.Proof.Proofs {
				cj := proof.Challenge
				computedC = addQ(computedC, &cj)

				vj := proof.ProofResponse
				wj := subQ(&vj, mulQ(big.NewInt(int64(j)), &cj))

				aj := mulP(powP(g, &vj), powP(alphaHat, &cj))
				bj := mulP(powP(k, wj), powP(betaHat, &cj))

				toBeHashed = append(toBeHashed, *aj)
				toBeHashed = append(toBeHashed, *bj)
				// TODO: Check A is already done in step5?

				checkB := isInRange(cj)
				checkC := isInRange(vj)

				if !checkB || !checkC {
					fmt.Println("we fucked the proof up")
				}
			}
			c := crypto.HMAC(constants.GetExtendedBaseHash(), 0x21, toBeHashed...)
			checkD := c.Cmp(computedC) == 0 // checking if c is computed correct

			if !checkD {
				fmt.Println("we fucked up d")
				fmt.Println(c.Text(16))
				fmt.Println(computedC.Text(16))
			}
		}
	}
}

func computeContestTotal(ballotSelections []models.BallotSelection) (*big.Int, *big.Int) {
	alphaHat := big.NewInt(1)
	betaHat := big.NewInt(1)

	for _, selection := range ballotSelections {
		alphaHat = mulP(alphaHat, &selection.Ciphertext.Pad)
		betaHat = mulP(betaHat, &selection.Ciphertext.Data)
	}

	return alphaHat, betaHat
}
