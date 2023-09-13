package gpio

import (
	"fmt"
	"log/slog"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
	"github.com/walnuts1018/machine-status-api/domain"
)

type client struct {
}

func NewClient() domain.GPIOClient {
	return client{}
}

func (c client) StartAlice() error {
	pwon, err := c.IsPwOn()
	if err != nil {
		return fmt.Errorf("failed to get read alice pw led voltage:%v", err)
	}
	if pwon {
		slog.Warn("alice is already running")
		return nil
	}

	pushTime := 800 * time.Millisecond
	err = c.pressPWSwitch(pushTime)
	if err != nil {
		return fmt.Errorf("failed to start alice: %w", err)
	}

	return nil
}

func (c client) StopAlice() error {
	pwon, err := c.IsPwOn()
	if err != nil {
		return fmt.Errorf("failed to get read alice pw led voltage:%v", err)
	}
	if !pwon {
		slog.Warn("alice is already stopped")
		return nil
	}

	pushTime := 800 * time.Millisecond

	err = c.pressPWSwitch(pushTime)
	if err != nil {
		return fmt.Errorf("failed to stop alice: %w", err)
	}

	timeout := time.After(5 * time.Minute)
	for {
		select {
		case <-timeout:
			slog.Warn("failed to shutdown alice, force stop")
			pushTime = 10 * time.Second
			err = c.pressPWSwitch(pushTime)
			if err != nil {
				return fmt.Errorf("failed to stop alice: %w", err)
			}
			return nil
		default:
			pwon, err := c.IsPwOn()
			if err != nil {
				return fmt.Errorf("failed to get read alice pw led voltage:%v", err)
			}
			if !pwon {
				return nil
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (c client) pressPWSwitch(pushTime time.Duration) error {
	err := rpio.Open()
	if err != nil {
		return fmt.Errorf("failed to open gpio: %w", err)
	}

	pwSwPin := rpio.Pin(21) // GPIO21, not pin number
	pwSwPin.Output()

	pwSwPin.High()
	time.Sleep(pushTime)
	pwSwPin.Low()
	pwSwPin.Input()

	err = rpio.Close()
	if err != nil {
		return fmt.Errorf("failed to close gpio: %w", err)
	}

	return nil
}

func (c client) IsPwOn() (bool, error) {
	err := rpio.Open()
	if err != nil {
		return false, fmt.Errorf("failed to open gpio: %w", err)
	}

	pwLedPin := rpio.Pin(16) // GPIO16, not pin number
	pwLedPin.Input()

	status := pwLedPin.Read()

	err = rpio.Close()
	if err != nil {
		return false, fmt.Errorf("failed to close gpio: %w", err)
	}

	return status == rpio.High, nil
}
