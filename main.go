package main

import "github.com/andrewsjg/go-loop/loopapi"

func main() {
	loopEn := loopapi.NewLoopEnergy()
	loopEn.Connect()
}
