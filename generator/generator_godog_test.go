package generator

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/cucumber/godog"
	"github.com/rs/zerolog/log"
)

func init() {
	config := GeneratorConfig{OpenAPIPath: "../examples/stores.yaml", OutputPath: "../build", ModuleName: "build"}
	err := GenerateServer(config)
	if err != nil {
		return
	}
	if _, err := os.Stat("../build"); err != nil {
		if os.IsNotExist(err) {
			log.Fatal()
		}
	}
	//cmd := exec.Command("exit")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = "go-open-api-generator/build/cmd"
	err = cmd.Run()
	if err != nil {
		log.Fatal()
	}
}

func iSendGETRequestTo(endpoint string) error {

	matcher, err := regexp.MatchString("/store/\\{id}", endpoint)
	if err != nil {
		panic(err)
	}
	url := "http://localhost:8000/" + endpoint
	data := strings.NewReader("")

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, data)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if endpoint == "/store" && resp.StatusCode != 404 {
		return fmt.Errorf("expected 404 received %d", resp.StatusCode)
	} else if matcher && resp.StatusCode != 501 {
		return fmt.Errorf("expected 501 received %d", resp.StatusCode)
	}

	return nil
}

func iSendPOSTRequestToWithPayload(endpoint, payload string) error {
	url := "http://localhost:8000/" + endpoint
	data := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, data)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 404 {
		return fmt.Errorf("expected 404 received %d", resp.StatusCode)
	}
	return nil
}

func iSendPUTRequestToWithPayload(endpoint, payload string) error {
	matcher, err := regexp.MatchString("/store/\\{id}", endpoint)
	if err != nil {
		panic(err)
	}

	url := "http://localhost:8000/" + endpoint
	data := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, data)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if endpoint == "/store" && resp.StatusCode != 404 {
		return fmt.Errorf("expected 404 received %d", resp.StatusCode)
	} else if matcher && resp.StatusCode != 501 {
		return fmt.Errorf("expected 501 received %d", resp.StatusCode)
	}
	return nil
}

func theResponseForUrlWithRequestMethodShouldBe(method, url string, statusCode int) error {
	matcher, err := regexp.MatchString("http://localhost:8000/store/\\{id}", url)
	if err != nil {
		panic(err)
	}
	if method == "GET" && url == "http://localhost:8000/store" {
		if statusCode != 404 {
			return fmt.Errorf("Expected 404 but received other status code")
		}
	} else if method == "GET" && matcher {
		if statusCode != 501 {
			return fmt.Errorf("Expected 404 but received other status code")
		}
	} else if method == "POST" {
		if statusCode != 404 {
			return fmt.Errorf("Expected 404 but received other status code")
		}
	} else if method == "PUT" {
		if statusCode != 501 {
			return fmt.Errorf("Expected 404 but received other status code")
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^I send GET request to "([^"]*)"$`, iSendGETRequestTo)
	ctx.Step(`^I send POST request to "([^"]*)" with payload "([^"]*)"$`, iSendPOSTRequestToWithPayload)
	ctx.Step(`^I send PUT request to "([^"]*)" with payload "([^"]*)"$`, iSendPUTRequestToWithPayload)
	ctx.Step(`^The response for url "([^"]*)" with request method "([^"]*)" should be (\d+)$`, theResponseForUrlWithRequestMethodShouldBe)
}
