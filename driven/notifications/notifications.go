package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Notifications structure
type Notifications struct {
	apiKey string
	host   string
}

// New creates new instance
func New(apiKey string, host string) Notifications {
	return Notifications{apiKey: apiKey, host: host}
}

// SendDataMsg sends data message
func (n *Notifications) SendDataMsg(topic string, data map[string]string) error {
	bodyJSON := make(map[string]interface{})
	bodyJSON["topic"] = topic
	bodyJSON["data"] = data

	bodyByteArr, err := json.Marshal(bodyJSON)
	if err != nil {
		log.Printf("notifications -> SendDataMsg: failed to Marshal body. Reason: %s", err.Error())
		return err
	}

	code, response, err := n.sendNotification(bytes.NewReader(bodyByteArr))
	if err != nil {
		log.Printf("notifications -> SendDataMsg: failed to send notification. Reason: %s", err.Error())
		return err
	}

	if *code == 200 {
		log.Println("notifications -> SendDataMsg: Success!")
		return nil
	}

	log.Printf("notifications -> SendDataMsg: request failed with code %d. Reason: %s", code, *response)
	return fmt.Errorf("%d: %s", code, *response)
}

// SendNotificationMsg sends notification message
func (n *Notifications) SendNotificationMsg(topic string, title string, body string, data map[string]string) error {
	bodyJSON := make(map[string]interface{})
	bodyJSON["topic"] = topic
	bodyJSON["subject"] = title
	bodyJSON["body"] = body
	bodyJSON["data"] = data

	bodyByteArr, err := json.Marshal(bodyJSON)
	if err != nil {
		log.Printf("notifications -> SendNotificationMsg: failed to Marshal body. Reason: %s", err.Error())
		return err
	}

	code, response, err := n.sendNotification(bytes.NewReader(bodyByteArr))
	if err != nil {
		log.Printf("notifications -> SendNotificationMsg: failed to send notification. Reason: %s", err.Error())
		return err
	}

	if *code == 200 {
		log.Println("notifications -> SendNotificationMsg: Success!")
		return nil
	}

	log.Printf("notifications -> SendNotificationMsg: request failed with code %d. Reason: %s", code, *response)
	return fmt.Errorf("%d: %s", code, *response)
}

func (n *Notifications) sendNotification(body io.Reader) (*int, *string, error) {
	if n.host == "" {
		return nil, nil, fmt.Errorf("missing host")
	}

	if n.apiKey == "" {
		return nil, nil, fmt.Errorf("missing internal api key")
	}

	if body == nil {
		return nil, nil, fmt.Errorf("missing body")
	}

	url := fmt.Sprintf("%s/notifications/api/int/message", n.host)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("INTERNAL-API-KEY", n.apiKey)
	client := &http.Client{Transport: &http.Transport{}}
	resp, err := client.Do(req)

	if err != nil {
		return nil, nil, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return nil, nil, err
	}

	responseString := string(bodyBytes)
	return &resp.StatusCode, &responseString, nil
}
