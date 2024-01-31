package models

type Ballot struct {
	ObjectId string
	Contests []BallotContest
}

type BallotContest struct {
	ObjectId         string
	CryptoHash       BigInt
	BallotSelections []BallotSelection
	Proof            RangeProof
}

type BallotSelection struct {
	ObjectId   string
	Ciphertext Ciphertext
	Proof      RangeProof
}

type RangeProof struct {
	Challenge  BigInt
	Proofs     []ChaumPedersenProof
	RangeLimit int
}

type ChaumPedersenProof struct {
	Challenge     BigInt
	ProofPad      BigInt
	ProofData     BigInt
	ProofResponse BigInt
}

type Ciphertext struct {
	Pad  BigInt
	Data BigInt
}
