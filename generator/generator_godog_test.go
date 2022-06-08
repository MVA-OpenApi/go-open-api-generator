package generator

import (
	"fmt"
	"github.com/cucumber/godog"
	"net/http"
	"strings"
)

func iRegisteredStoreWith(arg1 int) error {
	return godog.ErrPending
}

func iSendRequestTo(method, path string) error {
	data := strings.NewReader("id:456")

	// initialize http client
	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8000/store", data)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 405 {
		return fmt.Errorf("expected 405 received %d", resp.StatusCode)
	}
	return nil
}

func iSendRequestToWithPayload(method, payload, path string) error {
	data := strings.NewReader(payload)
	url := "http://localhost:8000/" + path

	// initialize http client
	client := &http.Client{}

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		panic(err)
	}

	// set the request header Content-Type for json
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 405 {
		return fmt.Errorf("expected 405 received %d", resp.StatusCode)
	}
	return nil
}

func theResponseShouldBe(response int, expected int) (bool, error) {
	if response != expected {
		return false, fmt.Errorf("expected response was %d, but actually got %d", expected, response)
	}
	return true, nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I registered store with (\d+)$`, iRegisteredStoreWith)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, iSendRequestTo)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)" with payload "([^"]*)"$`, iSendRequestToWithPayload)
	ctx.Step(`^The response should be (\d+)$`, theResponseShouldBe)
}
