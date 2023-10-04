.PHONY: send-event
send-event:
	go run ./cmd/cli/main.go -d '{"events":[{"id":"","createdAt":"2","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":"{\"count\":\"123\",\"line\":\"15\"}"},{"id":"","createdAt":"2","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":"{\"count\":\"234\",\"line\":\"27\"}"}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec"

.PHONY: start-mock
start-mock:
	go run ./cmd/mock/main.go