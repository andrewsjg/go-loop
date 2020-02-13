package main

import (
	"fmt"
	"os"
	"time"

	"github.com/andrewsjg/go-loop/loopapi"
)

func main() {

	// Get the serial and secret from the environment
	ELECSECRET := os.Getenv("LOOPELECSECRET")
	ELECSERIAL := os.Getenv("LOOPELECSERIAL")

	// No Gas Secret or serial provided. I dont have one to test with!
	loopEn := loopapi.NewLoopEnergy(ELECSERIAL, ELECSECRET, "", "", "www.your-loop.com", 443)
	loopEn.Connect()

	var lastElec float32

	for {
		if loopEn.Connected {
			time.Sleep(1 * time.Second)

			if loopEn.Electricty != lastElec {
				lastElec = loopEn.Electricty

				fmt.Println("Elec:", loopEn.Electricty)
				//loopEn.Disconnect()
				//break
			}
		}
	}

}
