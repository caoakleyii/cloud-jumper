package handler

import (
	"testing"
)

func TestShutdown(t *testing.T) {
	/*a := New()

	a.Get("/shutdown", Shutdown)

	server := &http.Server{
		Addr:    ":9999",
		Handler: a,
	}

	go func() {
		server.ListenAndServe()
	}()

	go func() {
		resp, err := http.Get("http://localhost:9999/shutdown")

		if err != nil {
			t.Errorf("TestShutdown errored when making request to test server \n\n %v", err)
		}

		if resp.StatusCode != 200 {
			t.Errorf("TestShutdown did not return an OK status. Returned: %v", resp.StatusCode)
		}
	}()

	signal.Notify(CSignal, os.Kill, os.Interrupt)
	<-CSignal
	close(CSignal)
	server.Shutdown(context.Background()) */

}
