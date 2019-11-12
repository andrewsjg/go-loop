package loopapi

import (
	"log"
	"os"
	"time"

	//socketio_client "github.com/zhouhui8915/go-socket.io-client"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
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

type connectMessage struct {
	serial		string	`json:"serial"`
	secret  	string  `json:"secret"`
	clientIP 	string  `json:"clientIp"`

}

// NewLoopEnergy - Initializes a new LoopEnergy object
func NewLoopEnergy() LoopEnergy {
	var theLoop LoopEnergy
	var cfg config

	// Setup the config
	// TODO: probably read these from a config file
	cfg.LoopPort = 443
	cfg.LoopServer = "www.your-loop.com"

	// Secret and Serial read from the environment
	cfg.Secret = os.Getenv("LOOPSECRET")
	cfg.Serial = os.Getenv("LOOPSERIAL")

	theLoop.Config = cfg

	return theLoop
}

// Connect - Connect to the LOOP API Endpoint
func (loopEn *LoopEnergy) Connect() bool {
	log.Println("Connecting")

	c, err := gosocketio.Dial(
		gosocketio.GetUrl(loopEn.Config.LoopServer, loopEn.Config.LoopPort, true),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("electric_realtime", func(h *gosocketio.Channel) {
		log.Println("Got Electric Data")
	})
	if err != nil {
		log.Fatal(err)
	}

	var msg connectMessage

	msg.clientIP = "127.0.0.1"
	msg.secret = loopEn.Config.Secret
	msg.serial = loopEn.Config.Serial

	c.Emit("subscribe_electric_realtime",msg)
	//client.Emit("subscribe_electric_realtime", "{ 'serial': "+loopEn.Config.Serial+",'clientIp': '127.0.0.1','secret': "+loopEn.Config.Secret+"}")

	time.Sleep(10 * time.Second)
	
	return true
}
