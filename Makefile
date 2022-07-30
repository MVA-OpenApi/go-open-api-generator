OUTPUT_NAME = "generator"
MODULE_NAME = "build"

generate:
	go run main.go generate $(OPEN_API_PATH)

generate-all-flags:
	go run main.go generate $(OPEN_API_PATH) -o . -n $(MODULE_NAME) -l -d

build:
	go build -o $(OUTPUT_NAME)

test:
	go test ./... -v