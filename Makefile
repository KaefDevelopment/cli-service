PROJECT_NAME=cli
BUILD_DIR=./bin
VERSION=$(git describe --tags --abbrev=0)

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
	go run ./cmd/cli/main.go -d '{"events":[{"id":"","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}},{"id":"","createdAt":"2","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"27"}}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec" -cli-version=true

.PHONY: start-mock
start-mock:
	go run ./cmd/mock/main.go

## unix - '{"events":[{"id":"","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}}'

## windows - "{\"events\":[{\"id\":\"\","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}}"


