[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.github.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Audit Logging Library</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/auditlog?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/auditlog)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/auditlog?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/auditlog/)


This is an audit logging library for [ContainerSSH](https://containerssh.github.io). Among others, it contains the encoder and decoder for the [ContainerSSH Audit Log Format](https://containerssh.github.io/audit/format/) written in Go. This readme will guide you through the process of using this library.

## Setting up a logging pipeline

This section will explain how to set up and use a logging pipeline. As a first step, you must create the logger. The easiest way to do that is to pass a config object:

```go
auditLogger, err := auditlog.New(cfg, logger)
```

The `cfg` variable must be of the type `auditlog.Config`. Here's an example configuration:

```go
config := auditlog.Config{
    Format:  "binary",
    Storage: "file",
    File: file.Config{
        Directory: "/tmp/auditlog",
    },
    Intercept: auditlog.InterceptConfig{
        Stdin:     true,
        Stdout:    true,
        Stderr:    true,
        Passwords: true,
    },
}
```
 
The `logger` variable must be an instance of `github.com/containerssh/log/logger`. The easiest way to create the logger is as follows:

```go
logger := logPipeline.NewLoggerPipeline(
    log.LevelDebug,
    ljson.NewLJsonLogFormatter(),
    os.Stdout,
)
```

Alternatively, you can also create the audit logger using the following factory method:

```go
logger := auditlog.NewLogger(
	intercept,
	encoder,
	storage,
	logger,
)
```

In this case `intercept` is of the type `InterceptConfig`, `encoder` is an instance of `codec.Encoder`, `storage` is an instance of `storage.WritableStorage`, and `logger` is the same logger as explained above. This allows you to create a custom pipeline.

### Writing to the pipeline

Once the audit logging pipeline is created you can then create your first entry for a new connection:

```go
connection, err := auditLogger.OnConnect(
    []byte("asdf"),
    net.TCPAddr{
        IP:   net.ParseIP("127.0.0.1"),
        Port: 2222,
        Zone: "",
    },
)
```

This will post a `connect` message to the audit log. The `connection` variable can then be used to send
subsequent connection-specific messages:

```go
connection.OnAuthPassword("foo", []byte("bar"))
connection.OnDisconnect()
```

The `OnNewChannelSuccess()` method also allows for the creation of a channel-specific audit logger that will log with the appropriate channel ID. 

## Retrieving and decoding messages

Once the messages are restored they can be retrieved by the same storage mechanism that was used to store them:

```go
storage, err := auditlog.NewStorage(config, logger)
if err != nil {
    log.Fatalf("%v", err)
}
// This only works if the storage type is not "none"
readableStorage := storage.(storage.ReadableStorage)
```

The readable storage will let you list audit log entries as well as fetch individual audit logs:

```go
logsChannel, errors := readableStorage.List()
for {
    finished := false
    select {
    case entry, ok := <-logsChannel:
        if !ok {
            finished = true
            break
        }
        // use entry.Name to reference a specific audit log
    case err, ok := <-errors:
        if !ok {
            finished = true
            break
        }
        if err != nil {
            // Handle err
        }
    }
    if finished {
        break
    }
}
```

Finally, you can fetch an individual audit log:

```go
reader, err := readableStorage.OpenReader(entry.Name)
if err != nil {
    // Handle error
}
```

The reader is now a standard `io.Reader`. 

## Decoding messages

Messages can be decoded with the reader as follows:

```go
// Set up the decoder
decoder := binary.NewDecoder()

// Decode messages
decodedMessageChannel, errorsChannel := decoder.Decode(reader)

for {
    finished := false
    select {
        // Fetch next message or error
        case msg, ok := <-decodedMessageChannel:
            if !ok {
                //Channel closed
                finished = true
                break
            } 
            //Handle messages
        case err := <-errorsChannel:
            if !ok {
                //Channel closed
                finished = true
                break
            } 
            // Handle error
    }
    if finished {
        break
    }
}
```

**Tip:** The `<-` signs are used with channels. They are used for async processing. If you are unfamiliar with them take a look at [Go by Example](https://gobyexample.com/channels).

**Note:** The Asciinema encoder doesn't have a decoder pair as the Asciinema format does not contain enough information to reconstruct the messages.

## Manually encoding messages

If you need to encode messages by hand without a logger pipeline you can do so with an encoder implementation. This is normally not needed. We have two encoder implementations: the binary and the Asciinema encoders. You can use them like this:

```go
encoder := binary.NewEncoder()
// Alternatively:
// encoder := asciinema.NewEncoder()

// Initialize message channel
messageChannel := make(chan message.Message)
// Initialize storage backend
storage := YourNewStorage()

go func() {
    err := encoder.Encode(messageChannel, storage)
    if err != nil {
        log.Fatalf("failed to encode messages (%w)", err)        
    }
}()

messageChannel <- message.Message{
    //Fill in message details here
}
close(messageChannel)
```

**Note:** The encoder will run until the message channel is closed, or a disconnect message is sent.

## Development

### Implementing an encoder and decoder

If you want to implement your own encoder for a custom format you can do so by implementing the `Encoder` interface in the [codec/abstract.go file](codec/abstract.go). Conversely, you can implement the `Decoder` interface to implement a decoder.

### Implementing a writable storage

In order to provide storages you must provide an `io.WriteCloser` with this added function:

```go
// Set metadata for the audit log. Can be called multiple times.
//
// startTime is the time when the connection started in unix timestamp
// sourceIp  is the IP address the user connected from
// username  is the username the user entered. The first time this method
//           is called the username will be nil, may be called subsequently
//           is the user authenticated.
SetMetadata(startTime int64, sourceIp string, username *string)
```

### Implementing a readable storage

In order to implement a readable storage you must implement the `ReadableStorage` interface in [storage/storage.go](storage/storage.go). You will need to implement the `OpenReader()` method to open a specific audit log and the `List()` method to list all available audit logs.