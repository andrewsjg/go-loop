package loopapi

import (
	"log"
	"os"

	socketio_client "github.com/zhouhui8915/go-socket.io-client"
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
func (loopEn *LoopEnergy) Connect() bool {
	log.Println("Connecting")

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}
	opts.Query["serial"] = loopEn.Config.Serial
	opts.Query["secret"] = loopEn.Config.Secret

	client, err := socketio_client.NewClient(loopEn.Config.LoopServer, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return false
	}

	client.On("electric_realtime", func() {
		log.Printf("on electric_realtime\n")
	})

	client.Emit("subscribe_gas_interval", "{ 'serial': "+loopEn.Config.Serial+",'clientIp': '127.0.0.1','secret': "+loopEn.Config.Secret+"}")

	for {
		log.Println("Listening")
	}
	return true
}
