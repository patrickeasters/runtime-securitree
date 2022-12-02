package decorator

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/patrickeasters/runtime-securitree/sds"
)

// LED colors by severity
var LEDSeverity = map[string][]int{
	"High":   {238, 99, 94},
	"Medium": {255, 169, 64},
	"Low":    {253, 216, 53},
	"Info":   {115, 161, 247},
}

var LEDDefault = []int{97, 242, 97}

type Decorator struct {
	LEDState    [][]int
	StripLength int
	Brightness  int
}

func NewDecorator(stripLen, brightness int) Decorator {
	state := make([][]int, stripLen)
	for i := range state {
		state[i] = LEDDefault
	}
	switch {
	case brightness > 255:
		brightness = 255
	case brightness < 0:
		brightness = 0
	}
	return Decorator{
		LEDState:    state,
		StripLength: stripLen,
		Brightness:  brightness,
	}
}

func (d *Decorator) PushLED(led []int) error {
	if len(d.LEDState) >= d.StripLength {
		d.LEDState = d.LEDState[1:]
	}
	d.LEDState = append(d.LEDState, led)
	// TODO: do API call to set state
	return nil
}

func (d *Decorator) MessageHandler(msg string) error {
	var event sds.PolicyEvent
	err := json.Unmarshal([]byte(msg), &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	led, ok := LEDSeverity[event.Severity.String()]
	if !ok {
		return errors.New("unknown event severity")
	}

	return d.PushLED(led)
}
