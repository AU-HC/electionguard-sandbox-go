package verification

import "fmt"

func verify(description string, ok bool) {
	if !ok {
		fmt.Printf("Verification error for check: %s", description)
	}
}
