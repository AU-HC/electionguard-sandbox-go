# ElectionGuard 2.0 Sandbox in Go
*Andreas Skriver Nielsen, Markus Valdemar Grønkjær Jensen, and Hans-Christian Kjeldsen*

## Overview
[ElectionGuard](https://github.com/microsoft/electionguard) is an open-source software development kit (SDK) by Microsoft,
that aims to improve the security and transparency of elections. The primary focus is that

- Individual voters can verify that their votes have been accurately recorded.
- Voters and observers can verify that all recorded votes have been accurately counted.

The aim of this repository is to build a sandbox environment to generate data for and verify range proofs not yet implemented in the hybrid 1.91.18 specification of ElectionGuard.

## Installation
As a prerequisite make sure to have installed Go, which can be downloaded [here](https://go.dev/doc/install). Afterwards download the verifier as a ZIP, or clone the repository from source:
```
$ git clone https://github.com/AU-HC/electionguard-sandbox-go.git 
```

## Usage
The sandbox is currently a command line utility tool, to generate and verify ballots the following command has to be executed.
```
$ go run main.go -p="path/to/manifest/"
```
It's important to note that the `-p` flag must be set, as it specifies the manifest path which specifies the contest, contest selection limit, selections, and selection limit.

The verifier also has alternate options which can be set, using the following flags:
- `-number` of type `int`: Specifies the amount of ballots generated and verified.
- `-output` of type `bool`: Specifies if the ballots should be saved as json in the `ballots/` directory.

For example to generate and verify 10 ballots using the `manifest.json` file, one of the following commands can be executed:
```
$ go run main.go -p="manifest.json" -number=10 
```
or (Windows)
```
$ go build main.go
$ electionguard-sandbox-go.exe -p="manifest.json" -number=10 
```
or (Mac/Linux)
```
$ go build main.go
$ ./electionguard-sandbox-go -p="manifest.json" -number=10 
```