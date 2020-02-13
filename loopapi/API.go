package loopapi

import (
	"sync"

	log "github.com/sirupsen/logrus"

	//socketio_client "github.com/zhouhui8915/go-socket.io-client"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type config struct {
	loopServer string
	loopPort   int
	elecSerial string
	elecSecret string
	gasSerial  string
	gasSecret  string
}

// LoopEnergy is the main container for data from Loop and the client
type LoopEnergy struct {
	config
	Connected bool
	client    *gosocketio.Client
	stop      chan bool

	Electricty float32
}

//RequestMessage is the data structure used sending the serial and secret along with the subscription requests
type RequestMessage struct {
	Serial   string `json:"serial"`
	Secret   string `json:"secret"`
	ClientIP string `json:"clientIp"`
}

//ElecDataMessage is the message returned from Loop
type ElecDataMessage struct {
	Lqi             int    `json:"lqi"`
	Rssi            int    `json:"rssi"`
	DeviceTimeStamp int    `json:"deviceTimestamp"`
	Inst            int    `json:"inst"`
	Serial          string `json:"serial"`
}

// NewLoopEnergy - Initializes a new LoopEnergy object
func NewLoopEnergy(elecSerial, elecSecret, gasSerial, gasSecret, loopServer string, loopPort int) LoopEnergy {
	var theLoop LoopEnergy

	// Setup the config
	theLoop.loopPort = loopPort
	theLoop.loopServer = loopServer

	theLoop.stop = make(chan bool)

	theLoop.elecSerial = elecSerial
	theLoop.elecSecret = elecSecret

	theLoop.gasSecret = gasSecret
	theLoop.gasSerial = gasSerial

	theLoop.Connected = false

	// Set logging
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.DebugLevel)

	return theLoop
}

// Connect - Connect to the LOOP API Endpoint
func (loopEn *LoopEnergy) connect(wg *sync.WaitGroup) {
	log.Info("Connecting")

	trans := transport.GetDefaultWebsocketTransport()

	c, err := gosocketio.Dial(gosocketio.GetUrl(loopEn.loopServer, loopEn.loopPort, true), trans)
	loopEn.client = c

	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Info("Connected to loop")

		// Subscribe to electricity readings
		msg := RequestMessage{loopEn.elecSerial, loopEn.elecSecret, "127.0.0.1"}
		emitErr := c.Emit("subscribe_electric_realtime", msg)

		// Change this!
		if emitErr != nil {
			log.Fatal(err)
		}

		// Subscribe to gas readings
		//msg := RequestMessage{loopEn.gasSerial, loopEn.gasSecret, "127.0.0.1"}
		//emitErr := c.Emit("subscribe_gas_interval", msg)

	})

	if err != nil {
		log.Fatal(err)
	}

	loopEn.Connected = true

	// Signal that the client is connected
	wg.Done()

	err = c.On("disconnect", func(h *gosocketio.Channel) {
		log.Info("Disconnected")
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

		log.Debug("Current Usage (kW):", currentUsage)

	})

	loopEn.client = c
	loopEn.Connected = true

	// Block on the signal channel

loop:
	for {
		select {
		case <-loopEn.stop: // triggered when the stop channel is closed
			log.Debug("Disconnecting")
			break loop // exit
		}
	}
}

//Connect connects to Loop and subscribes to the channels for data
func (loopEn *LoopEnergy) Connect() bool {

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go loopEn.connect(wg)
	wg.Wait()

	return true
}

// Disconnect from loop. Terminates the client and send a message to the stop channel to
// terminate the goroutine running the client
func (loopEn *LoopEnergy) Disconnect() {
	loopEn.client.Close()
	loopEn.stop <- true
}
