package support

import (
	"fmt"
	"net/http"
	"time"
)

func WaitForReady(timeout time.Duration, endpoint string) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequest(http.MethodGet, endpoint, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}

		if time.Since(startTime) >= timeout {
			return fmt.Errorf("timeout when waiting for %s", endpoint)
		}

		time.Sleep(250 * time.Millisecond)
	}
}
