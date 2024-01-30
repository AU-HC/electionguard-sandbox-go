package models

import "math/big"

type Ballot struct {
	ObjectId string
	Code     string
	Contests []BallotContest
}

type BallotContest struct {
	ObjectId               string
	CryptoHash             big.Int
	BallotSelections       []BallotSelection
	Proof                  RangeProof
	CiphertextAccumulation Ciphertext
}

type BallotSelection struct {
	ObjectId   string
	Ciphertext Ciphertext
	Proof      RangeProof
}

type RangeProof struct {
	Challenge  big.Int
	Proofs     []ChaumPedersenProof
	RangeLimit int
}

type ChaumPedersenProof struct {
	Challenge     big.Int
	ProofPad      big.Int
	ProofData     big.Int
	ProofResponse big.Int
}

type Ciphertext struct {
	Pad  big.Int
	Data big.Int
}
