package vikunja

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
)

// ConsumeWebhookCallback consumes a webhook callback and calls the callback function
func ConsumeWebhookCallback(body io.ReadCloser, callback func(webhook WebhookCallback) error) error {
	defer body.Close()
	bytes, err := io.ReadAll(body)
	if err != nil {
		return err
	}

	// Figure out if it has body and i need to unpack it or its already unpacked
	obj := map[string]any(nil)
	if err := json.Unmarshal(bytes, &obj); err != nil {
		return err
	}

	event := WebhookCallback{}
	if b, ok := obj["body"]; ok {
		bStr := b.(string)

		if err := json.Unmarshal([]byte(bStr), &event); err != nil {
			return err
		}
	} else {
		if err := json.Unmarshal(bytes, &event); err != nil {
			return err
		}
	}

	if err := callback(event); err != nil {
		return err
	}

	return nil
}

// RegisterVikunjaWebhookHandler registers a webhook handler for Vikunja Webhook
/*
Typical usage is something like:
l := utils.GetInitLogger()

if err := vikunja.RegisterVikunjaWebhookHandler("/", SomeWebhookHandler); err != nil {
	log.Fatal(err)
}

utils.RunAPIServer(8080)
*/
func RegisterVikunjaWebhookHandler(path string, callback func(Webhook WebhookCallback, c *Client) error) error {
	l := slog.Default().With("path", path)
	l.Info("Registering vikunja webhook handler")

	c, err := GetVikunjaAPIClient("", "")
	if err != nil {
		return err
	}

	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := ConsumeWebhookCallback(r.Body, func(event WebhookCallback) error { return callback(event, c) })
			if err != nil {
				l.Error(err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	return nil
}
