package verification

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	"fmt"
	"math/big"
)

func VerifyStep5(ballots []models.Ballot, publicKey crypto.PublicKey) {
	//g := publicKey.G
	//k := publicKey.K

	for _, ballot := range ballots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				alpha := selection.Ciphertext.Pad
				beta := selection.Ciphertext.Data

				checkAOne := isValidResidue(alpha)
				checkATwo := isValidResidue(beta)

				c := crypto.HMAC(constants.GetExtendedBaseHash(), 0x21, publicKey.K, selection.Ciphertext.Pad, selection.Ciphertext.Data)
				computedC := big.NewInt(0)
				for _, proof := range selection.Proof.Proofs {
					cj := proof.Challenge
					computedC = addQ(computedC, &cj)

					vj := proof.ProofResponse
					// wj := subQ(&vj, mulQ(big.NewInt(int64(j)), &cj))

					// aj := mulP(powP(g, &vj), powP(&alpha, &cj))
					// bj := mulP(powP(k, wj), powP(&beta, &cj))

					checkB := isInRange(cj)
					checkC := isInRange(vj)

					if !checkB || !checkC {
						fmt.Println("we fucked the proof up")
					}
				}

				checkD := c.Cmp(computedC) == 0 // checking if c is computed correct

				if !checkAOne {
					fmt.Println("we fucked up a 1")
				}

				if !checkATwo {
					fmt.Println("we fucked up a 2")
				}

				if !checkD {
					fmt.Println("we fucked up d")
					fmt.Println(c.Text(16))
					fmt.Println(computedC.Text(16))
				}
			}
		}
	}
}
