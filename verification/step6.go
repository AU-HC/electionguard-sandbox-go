package verification

import (
	"electionguard-sandbox-go/constants"
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/models"
	mod "electionguard-sandbox-go/modular_arithmetic"
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
				computedC = mod.AddQ(computedC, &cj)

				vj := proof.ProofResponse
				wj := mod.SubQ(&vj, mod.MulQ(big.NewInt(int64(j)), &cj))

				aj := mod.MulP(mod.PowP(g, &vj), mod.PowP(alphaHat, &cj))
				bj := mod.MulP(mod.PowP(k, wj), mod.PowP(betaHat, &cj))

				toBeHashed = append(toBeHashed, *aj)
				toBeHashed = append(toBeHashed, *bj)
				// TODO: Check A is already done in step5?

				verify("(6.B) ...", mod.IsInRange(cj))
				verify("(6.C) ...", mod.IsInRange(vj))
			}
			c := crypto.HMAC(constants.GetExtendedBaseHash(), 0x21, toBeHashed...)
			verify("(5.D) ...", c.Cmp(computedC) == 0)
		}
	}
}

func computeContestTotal(ballotSelections []models.BallotSelection) (*big.Int, *big.Int) {
	alphaHat := big.NewInt(1)
	betaHat := big.NewInt(1)

	for _, selection := range ballotSelections {
		alphaHat = mod.MulP(alphaHat, &selection.Ciphertext.Pad)
		betaHat = mod.MulP(betaHat, &selection.Ciphertext.Data)
	}

	return alphaHat, betaHat
}
