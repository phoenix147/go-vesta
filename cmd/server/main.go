package main

import (
	"log/slog"
	"os"

	vapi "github.com/phoenix147/go-vesta"
)

func main() {

	debug := os.Getenv("DEBUG")

	if debug == "1" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	vncAddress := os.Getenv("VNC_ADDR")
	vncPassword := os.Getenv("VNC_PASSWORD")
	httpAddr := os.Getenv("HTTP_ADDR")

	vestatherm := vapi.NewVestatherm(vncAddress, vncPassword)

	slog.Debug("connecting to vestatherm", "address", vncAddress)
	err := vestatherm.Connect()
	if err != nil {
		slog.Error("couldn't connect to vestatherm", "err", err)
	}

	slog.Debug("serving http api", "addr", httpAddr)
	vapi.RunApi(httpAddr, &vestatherm)
}
