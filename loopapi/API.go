package loopapi

import (
	log "github.com/sirupsen/logrus"

	//socketio_client "github.com/zhouhui8915/go-socket.io-client"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type config struct {
	loopServer string
	loopPort   int
	serial     string
	secret     string
}

type LoopEnergy struct {
	config
	Connected bool
	Client    *gosocketio.Client
	stop      chan bool

	Electricty float32
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
func NewLoopEnergy(serial, secret, loopServer string, loopPort int) LoopEnergy {
	var theLoop LoopEnergy

	// Setup the config
	theLoop.loopPort = loopPort
	theLoop.loopServer = loopServer

	theLoop.stop = make(chan bool)

	theLoop.serial = serial
	theLoop.secret = secret

	theLoop.Connected = false

	// Set logging
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	log.SetLevel(log.InfoLevel)

	return theLoop
}

// Connect - Connect to the LOOP API Endpoint
func (loopEn *LoopEnergy) connect() {
	log.Info("Connecting")

	c, err := gosocketio.Dial(
		gosocketio.GetUrl(loopEn.loopServer, loopEn.loopPort, true),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Info("Connected to loop")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("disconnect", func(h *gosocketio.Channel) {
		log.Println("Disconnected")
		loopEn.Connected = false
		loopEn.stop <- true

	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("electric_realtime", func(h *gosocketio.Channel, args ElecDataMessage) {
		log.Debug("Got Electric Data")
		//log.Println(args)

		currentUsage := float32(args.Inst) / 1000.0
		loopEn.Electricty = currentUsage

		log.Debug("Current Usage (kW): ", currentUsage)

	})
	if err != nil {
		log.Fatal(err)
	}

	/*
		time.Sleep(1 * time.Second)
		go subElec(c, loopEn.Config.Secret, loopEn.Config.Serial)

		time.Sleep(30 * time.Second) */

	loopEn.Client = c
	loopEn.Connected = true

	var msg RequestMessage
	msg.ClientIP = "127.0.0.1"
	msg.Secret = loopEn.secret
	msg.Serial = loopEn.serial

	err = c.Emit("subscribe_electric_realtime", msg)

	// Block on the signal channel
loop:

	for {
		select {
		case <-loopEn.stop: // triggered when the stop channel is closed
			break loop // exit
		}
	}
}

func (loopEn *LoopEnergy) Connect() bool {
	go loopEn.connect()
	return true
}

func (loopEn *LoopEnergy) Disconnect() {
	loopEn.stop <- true
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
	log.Debug("Waiting for Data")
}
