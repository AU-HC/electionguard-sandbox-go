package main

import (
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/generation"
	"electionguard-sandbox-go/verification"
	"flag"
	"fmt"
)

func main() {
	// Getting path to manifest file
	pathPtr := flag.String("path", "", "path to manifest file used to generate ballots.")
	amountOfBallotsPtr := flag.Int("number", 1, "amount of ballots to generate and verify.")
	flag.Parse()

	// Loading manifest and generating ballots based on manifest
	manifest := generation.LoadManifest(*pathPtr)
	electionPublicKey, _ := crypto.GenerateKeyPair()
	ballots := generation.GenerateBallots(manifest, *amountOfBallotsPtr, electionPublicKey)

	// Sandbox "verification" of ballots
	verification.VerifyStep5(ballots, electionPublicKey)
	verification.VerifyStep6(ballots, electionPublicKey)

	fmt.Println("We will see you in person.")
}
