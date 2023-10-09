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

### How to send data

-d string                                                                     
Request data in JSON format string

-k string                                                           
authorization key (plugin_id)

-s string                                                                     
Http address for sending events (default "http://localhost:8181/events")

### Example:

[your binary file] -d '{"events":[{"id":"","createdAt":"2","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"12"}},{"id":"","createdAt":"2","type":"2","project":"2","projectBaseDir":"/mnt/c/Users/jaros/GolandProjects/tts","language":"golang","target":"2","branch":"","timezone":"2","params":{"count":"17"}}]}' -k "346d7f75-4b20-4166-8577-e656cdf3caea" -s "https://kaef.io/api/plugin/v1/events"