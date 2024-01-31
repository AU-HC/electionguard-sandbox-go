package output_strategy

import (
	"electionguard-sandbox-go/models"
	"encoding/json"
	"fmt"
	"os"
)

type Strategy interface {
	OutputBallots(ballots []models.Ballot)
}

// NoOutput Outputs nothing and
type NoOutput struct {
}

func (s NoOutput) OutputBallots(ballots []models.Ballot) {

}

// FolderOutput Outputs and folder with json files of ballots
type FolderOutput struct {
	xd string
}

func (s FolderOutput) OutputBallots(ballots []models.Ballot) {
	err := os.RemoveAll("ballots/")
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("ballots/", 0755)
	if err != nil {
		panic(err)
	}

	for i, ballot := range ballots {
		jsonBytes, err := json.MarshalIndent(ballot, "", "  ")
		if err != nil {
			panic(err)
		}

		outputPath := fmt.Sprintf("ballots/ballot%d.json", i)

		err = os.WriteFile(outputPath, jsonBytes, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func MakeOutputStrategy(output bool) Strategy {
	if output {
		return FolderOutput{xd: "e"}
	}
	return NoOutput{}
}
