package vesta

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/kward/go-vnc"
	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
)

type Vestatherm struct {
	vncAddr string
	vncPass string
	vncConn *vnc.ClientConn
	vncMsg  chan vnc.ServerMessage
}

func NewVestatherm(vncAddr string, vncPass string) Vestatherm {
	return Vestatherm{
		vncAddr: vncAddr,
		vncPass: vncPass,
	}
}

func (v *Vestatherm) GetConnection() *vnc.ClientConn {
	return v.vncConn
}

func (v *Vestatherm) GetServerMessageChan() chan vnc.ServerMessage {
	return v.vncMsg
}

func (v *Vestatherm) Connect() error {
	// Establish TCP connection to VNC server.
	nc, err := net.Dial("tcp", v.vncAddr)
	if err != nil {
		return fmt.Errorf("Error connecting to VNC host. %w", err)
	}

	// Negotiate connection with the server.
	vcc := vnc.NewClientConfig(v.vncPass)
	vcc.ServerMessageCh = make(chan vnc.ServerMessage)
	vc, err := vnc.Connect(context.Background(), nc, vcc)
	if err != nil {
		return fmt.Errorf("Error negotiating connection to VNC host. %w", err)
	}

	pixelFormat := vnc.PixelFormat{
		BPP:        32,
		Depth:      24,
		BigEndian:  false,
		TrueColor:  true,
		RedMax:     255,
		GreenMax:   255,
		BlueMax:    255,
		RedShift:   16,
		GreenShift: 8,
		BlueShift:  0,
	}

	err = vc.SetPixelFormat(pixelFormat)
	if err != nil {
		return fmt.Errorf("Could not set pixel format: %w", err)
	}

	v.vncConn = vc
	v.vncMsg = vcc.ServerMessageCh

	return nil
}

func (v *Vestatherm) click(x, y int) error {
	err := v.vncConn.PointerEvent(buttons.Left, uint16(x), uint16(y))
	if err != nil {
		return fmt.Errorf("Error pressing button %w", err)
	}
	err = v.vncConn.PointerEvent(buttons.None, uint16(x), uint16(y))
	if err != nil {
		return fmt.Errorf("Error releasing button %w", err)
	}

	return nil
}

func (v *Vestatherm) keyPress(key keys.Key) error {
	err := v.vncConn.KeyEvent(key, vnc.PressKey)
	if err != nil {
		return fmt.Errorf("Error pressing key %w", err)
	}
	err = v.vncConn.KeyEvent(key, vnc.ReleaseKey)
	if err != nil {
		return fmt.Errorf("Error releasing key %w", err)
	}

	return nil
}

func (v *Vestatherm) GoToHome() error {
	slog.Debug("Changing screen to home")
	return v.click(60, 30)
}

func (v *Vestatherm) GoToTiledStove() error {
	slog.Debug("Changing screen to tiled stove")
	return v.click(170, 30)
}

func (v *Vestatherm) GoToHeatpump() error {
	slog.Debug("Changing screen to heatpump")
	return v.click(300, 30)
}

func (v *Vestatherm) GoToTrend() error {
	slog.Debug("Changing screen to trend")
	return v.click(400, 30)
}

func (v *Vestatherm) GoToNotifications() error {
	slog.Debug("Changing screen to notifications")
	return v.click(500, 30)
}

func (v *Vestatherm) GoToSettings() error {
	slog.Debug("Changing screen to settings")
	return v.click(650, 30)
}

func (v *Vestatherm) GoToBuffer() error {
	slog.Debug("Changing screen to buffer")
	err := v.GoToHeatpump()
	if err != nil {
		return err
	}

	return v.click(345, 420)
}

func (v *Vestatherm) StartDefrosting() error {
	slog.Debug("Start defrosting")
	err := v.GoToHome()
	if err != nil {
		return err
	}

	return v.click(80, 100)
}

func (v *Vestatherm) SetWaterTemperature(temp int) error {
	slog.Debug("Setting water temperature to ", "temp", temp)

	err := v.GoToBuffer()
	if err != nil {
		return err
	}

	// Click on the water temperature input field
	err = v.click(510, 80)
	if err != nil {
		return err
	}

	// Enter the new temperature digit by digit
	for i, d := range strconv.Itoa(temp) {
		if i > 0 {
			time.Sleep(100 * time.Millisecond)
		}
		if err := v.keyPress(keys.Digit0 + keys.Key(d-'0')); err != nil {
			return err
		}
	}

	// Confirm the new water temperature by clicking the ENTER button
	err = v.click(475, 400)

	return err

}
