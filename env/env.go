package env

// Environment management

import (
	"os"
)

type Environment = string

const (
	// Disconnected from other systemss
	TST Environment = "TST"
	// Connected to SendGrid API, but not to database
	GTU = "GTU"
	// Production
	PRD = "PRD"
)

var ENV Environment = initEnv()

func initEnv() Environment {
	env := os.Getenv("CLI_ENV")
	if len(env) == 0 {
		// PRD is deault environment
		env = PRD
	}
	return env
}

// Returns true if the CLI is ran in a test environment
func IsTest() bool {
	return ENV != PRD
}

func IsPRD() bool {
	return ENV == PRD
}

func IsGTU() bool {
	return ENV == GTU
}

func IsTST() bool {
	return ENV == TST
}
