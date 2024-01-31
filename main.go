package main

import (
	"electionguard-sandbox-go/crypto"
	"electionguard-sandbox-go/generation"
	"electionguard-sandbox-go/output_strategy"
	"electionguard-sandbox-go/verification"
	"flag"
	"fmt"
	"os"
)

func main() {
	// Getting path to manifest file
	pathPtr := flag.String("path", "", "path to manifest file used to generate ballots.")
	amountOfBallotsPtr := flag.Int("number", 1, "amount of ballots to generate and verify.")
	outputPtr := flag.Bool("output", false, "output generated ballots to folder")
	flag.Parse()

	// Path to manifest should be set
	if *pathPtr == "" {
		os.Exit(1)
	}

	// Loading manifest and generating ballots based on manifest
	manifest := generation.LoadManifest(*pathPtr)
	electionPublicKey, _ := crypto.GenerateKeyPair()
	ballots := generation.GenerateBallots(manifest, *amountOfBallotsPtr, electionPublicKey)

	// Create output strategy based on flag and output ballots
	outputStrategy := output_strategy.MakeOutputStrategy(*outputPtr)
	outputStrategy.OutputBallots(ballots)

	// Sandbox "verification" of ballots
	verification.VerifyStep5(ballots, electionPublicKey)
	verification.VerifyStep6(ballots, electionPublicKey)

	fmt.Println("Generation and verification of ballots finished successfully.")
}
