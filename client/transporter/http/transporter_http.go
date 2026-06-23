package http_transporter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Liphium/neoroute/client"
)

type HTTPTransporter struct {
	sender    client.Sender
	sendMutex sync.Mutex
}

func NewHTTPTransporter(s client.Sender, method string, u *url.URL) *HTTPTransporter {

	t := &HTTPTransporter{
		sender: s,
	}

	t.sender.SetSendFunc(func(data []byte) error {
		t.sendMutex.Lock()
		defer t.sendMutex.Unlock()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewReader(data))
		if err != nil {
			return err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		// Check for transporter errors
		if resp.StatusCode != http.StatusOK {
			return errors.New("received non ok status " + resp.Status + ": " + string(bodyBytes))
		}

		// Let sender handle the response routing
		go t.sender.Handle(bodyBytes)

		return nil
	})

	return t
}
