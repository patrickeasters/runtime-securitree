package wled

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	BaseURL    string
}

func NewClient(baseURL string) Client {
	return Client{
		httpClient: &http.Client{
			Transport: http.DefaultTransport.(*http.Transport).Clone(),
		},
		BaseURL: baseURL,
	}
}

// Just a subset of State attributes needed for this project
// See docs for more: https://kno.wled.ge/interfaces/json-api/
type State struct {
	On         *bool     `json:"on,omitempty"`
	Brightness *int      `json:"bri,omitempty"`
	Transition *int      `json:"transition,omitempty"`
	Preset     *int      `json:"ps,omitempty"`
	Playlist   *int      `json:"pl,omitempty"`
	Segment    []Segment `json:"seg,omitempty"`
}

type Segment struct {
	Individual [][]int `json:"i,omitempty"`
}

func (c Client) SetLEDs(brightness int, leds [][]int) error {
	on := true
	body := State{
		On:         &on,
		Brightness: &brightness,
		Segment:    []Segment{{Individual: leds}},
	}
	postBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal WLED post body: %w", err)
	}
	resp, err := c.httpClient.Post(c.BaseURL+"/json/state", "application/json", bytes.NewReader(postBody))
	if err != nil {
		return fmt.Errorf("failed to set LED state: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received invalid status when setting LED state: %d", resp.StatusCode)

	}
	return nil
}
