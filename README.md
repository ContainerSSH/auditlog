[![ContainerSSH - Launch Containers on Demand](https://containerssh.github.io/images/logo-for-embedding.svg)](https://containerssh.github.io/)

<!--suppress HtmlDeprecatedAttribute -->
<h1 align="center">ContainerSSH Audit Log Encoder and Decoder Library</h1>

[![Go Report Card](https://goreportcard.com/badge/github.com/containerssh/auditlog?style=for-the-badge)](https://goreportcard.com/report/github.com/containerssh/auditlog)
[![LGTM Alerts](https://img.shields.io/lgtm/alerts/github/ContainerSSH/auditlog?style=for-the-badge)](https://lgtm.com/projects/g/ContainerSSH/auditlog/)


This is an encoder and decoder library for the [ContainerSSH Audit Log Format](https://containerssh.github.io/audit/format/) written in Go. In order to use it you will need depend on `github.com/containerssh/auditlog`.

## Encoding messages

Messages can be encoded with a format encoder, for example:

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
        log.Fatalf("failed to encode messages (%v)", err)        
    }
}()

messageChannel <- &message.Message{
    //Fill in message details here
}
close(messageChannel)
```

**Note:** The encoder will run until the message channel is closed, or a disconnect message is sent.

### Implementing a storage

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

## Decoding messages

Messages can be decoded as follows:

```go
// Initialize the reader
// reader := set up your io.Reader here

// Set up the decoder
decoder := binary.NewDecoder()

// Decode messages
decodedMessageChannel, errorsChannel, done := decoder.Decode(pipeReader)

for {
    select {
        // Fetch next message or error
        case msg := <-decodedMessageChannel:
            //Handle messages
        case err := <-errorsChannel:
            // Handle error
    }
    select {
        case <- done:
            // Break cycle
            break
        default: 
            // Continue cycle
    }
}
```

**Note:** The Asciinema encoder doesn't have a decoder pair as the Asciinema format does not contain enough information to reconstruct the messages.
