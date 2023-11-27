# CLI-service

## How to debug locally

Start remote server mock:

```shell
make start-mock
```

Send test event through cli

```shell
make send-event
```

## How to get cli version

[your binary file] version

## How to send event

Usage:                                                                                                                                      
event [flags]

Flags:                                                                                                                                      
-k, --auth-key string     Authorization key                                                                                              
-d, --data string          Request data in JSON format string                                                                             
-h, --help                 help for cli-event                                                                                             
-s, --server-host string   Http address for sending events (default "https://nautime.io/api/plugin/v1/events?source=cli&version=$version")

### Example

[your binary file] event -d '{"events":[{"id":"","createdAt":"3","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caec"

