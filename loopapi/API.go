package loopapi

import (
	"log"
	"os"
)

type config struct {
	LoopServer string
	LoopPort   int
	Serial     string
	Secret     string
}

type LoopEnergy struct {
	Config config
}

// NewLoopEnergy - Initializes a new LoopEnergy object
func NewLoopEnergy() LoopEnergy {
	var theLoop LoopEnergy
	var cfg config

	// Setup the config
	// TODO: probably read these from a config file
	cfg.LoopPort = 443
	cfg.LoopServer = "https://www.your-loop.com"

	// Secret and Serial read from the environment
	cfg.Secret = os.Getenv("LOOPSECRET")
	cfg.Serial = os.Getenv("LOOPSERIAL")

	theLoop.Config = cfg

	return theLoop
}

// Connect - Connect to the LOOP API Endpoint
func (loopEng *LoopEnergy) Connect() bool {
	log.Println("TEST")

	return true
}
