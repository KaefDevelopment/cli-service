PROJECT_NAME=cli
BUILD_DIR=./bin
VERSION=$(shell git describe --tags --abbrev=0)

# go tool dist list
WINDOWS=windows/386 windows/amd64 windows/arm
DARWIN=darwin/amd64 darwin/arm64
LINUX=linux/386 linux/amd64 linux/arm linux/arm64
PLATFORMS=$(WINDOWS) $(LINUX) $(DARWIN)

run: build-all

.PHONY: build-all
build-all: $(PLATFORMS)

$(WINDOWS): EXT=.exe
$(PLATFORMS): split=$(subst /, ,$@)
$(PLATFORMS): OS=$(word 1,$(split))
$(PLATFORMS): ARCH=$(word 2,$(split))
$(PLATFORMS): ARTIFACT_NAME=$(PROJECT_NAME)-$(OS)-$(ARCH)$(EXT)
$(PLATFORMS):
	env GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=1 go build -ldflags="-X main.version=$(VERSION)" -o $(BUILD_DIR)/$(ARTIFACT_NAME) cmd/cli/main.go

.PHONY: zip-artifacts
zip-artifacts: $(foreach f,$(wildcard $(BUILD_DIR)/*[^zip]),$(f).zip)

$(BUILD_DIR)/%.zip:
	@cd $(BUILD_DIR) && zip $*.zip $*

.PHONY: send-event
send-event:
	go run ./cmd/cli/main.go event -d '{"events":[{"pluginId":"346d7f75-4b20-4166-8577-e656cdf3caec","id":"","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec" -s "http://localhost:8181/events"

.PHONY: send-empty
send-empty:
	go run ./cmd/cli/main.go event -d '{"events":[]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec" -s "http://localhost:8181/events"

.PHONY: send-bad-event
send-bad-event:
	go run ./cmd/cli/main.go event -d '{"events":[{"id":"","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12}}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec" -s "http://localhost:8181/events"

.PHONY: send-not-authorized
send-not-authorized:
	go run ./cmd/cli/main.go event -d '{"events":[{"id":"","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec2" -s "http://localhost:8181/events" -a=false

.PHONY: version
version:
	go run ./cmd/cli/main.go version

.PHONY: start-mock
start-mock:
	go run ./cmd/mock/main.go

.PHONY: test
test:
	@go test ./... -race -v

.PHONY: generate
generate:
	go generate ./...
