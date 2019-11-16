package loopapi

import (
	"log"
	"os"
	"time"

	//socketio_client "github.com/zhouhui8915/go-socket.io-client"
	gosocketio "github.com/graarh/golang-socketio"
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

type RequestMessage struct {
	Serial   string `json:"serial"`
	Secret   string `json:"secret"`
	ClientIP string `json:"clientIp"`
}

//{"lqi":43,"rssi":-92,"deviceTimestamp":1573945716,"inst":1660,"serial":"53e8"}
type ElecDataMessage struct {
	Lqi             int    `json:"lqi"`
	Rssi            int    `json:"rssi"`
	DeviceTimeStamp int    `json:"deviceTimestamp"`
	Inst            int    `json:"inst"`
	Serial          string `json:"serial"`
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
		log.Println("Connected to loop")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("disconnect", func(h *gosocketio.Channel) {
		log.Println("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("electric_realtime", func(h *gosocketio.Channel, args ElecDataMessage) {
		log.Println("Got Electric Data")
		//log.Println(args)

		currentUsage := float32(args.Inst) / 1000.0
		log.Printf("Current Usage (kW): %.3f ", currentUsage)

	})
	if err != nil {
		log.Fatal(err)
	}

	/*
		var msg connectMessage

		msg.clientIP = "127.0.0.1"
		msg.secret = loopEn.Config.Secret
		msg.serial = loopEn.Config.Serial

		c.Emit("subscribe_electric_realtime",msg)
	*/

	time.Sleep(1 * time.Second)
	go subElec(c, loopEn.Config.Secret, loopEn.Config.Serial)

	time.Sleep(30 * time.Second)

	return true
}

func subElec(c *gosocketio.Client, secret string, serial string) {
	var msg RequestMessage
	msg.ClientIP = "127.0.0.1"
	msg.Secret = secret
	msg.Serial = serial

	//result,err := c.Ack("subscribe_electric_realtime",msg, 5 * time.Second)
	err := c.Emit("subscribe_electric_realtime", msg)

	if err != nil {
		log.Fatal("Emit ERRROR: " + err.Error())
	}
	log.Println("Waiting for Data")
}
