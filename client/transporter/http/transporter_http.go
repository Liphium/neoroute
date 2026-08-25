package http

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

// ApplyHTTP makes a sender send Neoroute requests over HTTP, using the given method and URL.
func ApplyHTTP(s client.Sender, method string, u *url.URL) {
	sendMutex := sync.Mutex{}

	s.SetSendFunc(func(data []byte) error {
		sendMutex.Lock()
		defer sendMutex.Unlock()
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
		go s.Handle(bodyBytes)

		return nil
	})
}
