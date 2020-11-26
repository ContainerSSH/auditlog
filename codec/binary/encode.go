package binary

import (
	"compress/gzip"
	"fmt"
	"net"

	"github.com/containerssh/geoip"

	"github.com/containerssh/auditlog/storage"

	"github.com/containerssh/auditlog/codec"
	"github.com/containerssh/auditlog/message"

	"github.com/fxamacker/cbor"
)

// NewEncoder creates an encoder that encodes messages in CBOR+GZIP format as documented
//            on https://containerssh.github.io/advanced/audit/format/
func NewEncoder(geoIPProvider geoip.LookupProvider) codec.Encoder {
	return &encoder{
		geoIPProvider: geoIPProvider,
	}
}

type encoder struct {
	geoIPProvider geoip.LookupProvider
}

func (e *encoder) GetMimeType() string {
	return "application/octet-stream"
}

func (e *encoder) GetFileExtension() string {
	return ""
}

func (e *encoder) Encode(messages <-chan message.Message, storage storage.Writer) error {
	var gzipHandle *gzip.Writer
	var encoder *cbor.Encoder
	gzipHandle = gzip.NewWriter(storage)
	encoder = cbor.NewEncoder(gzipHandle, cbor.EncOptions{})
	if err := encoder.StartIndefiniteArray(); err != nil {
		return fmt.Errorf("failed to start infinite array (%w)", err)
	}

	startTime := int64(0)
	var ip = ""
	var country = "XX"
	var username *string
	for {
		msg, ok := <-messages
		if !ok {
			break
		}
		if startTime == 0 {
			startTime = msg.Timestamp
		}
		e.storeMetadata(msg, storage, &ip, &country, &startTime, username)
		if err := encoder.Encode(&msg); err != nil {
			return fmt.Errorf("failed to encode audit log message (%w)", err)
		}
		if msg.MessageType == message.TypeDisconnect {
			break
		}
	}
	if err := encoder.EndIndefinite(); err != nil {
		return fmt.Errorf("failed to end audit log infinite array (%w)", err)
	}
	if err := gzipHandle.Flush(); err != nil {
		return fmt.Errorf("failed to flush audit log gzip stream (%w)", err)
	}
	if err := storage.Close(); err != nil {
		return fmt.Errorf("failed to close audit log gzip stream (%w)", err)
	}
	return nil
}

func (e *encoder) storeMetadata(
	msg message.Message,
	storage storage.Writer,
	ip *string,
	country *string,
	startTime *int64,
	username *string,
) {
	switch msg.MessageType {
	case message.TypeConnect:
		remoteAddr := msg.Payload.(message.PayloadConnect).RemoteAddr
		ip = &remoteAddr
		isoCountry := e.geoIPProvider.Lookup(net.ParseIP(*ip))
		country = &isoCountry
		storage.SetMetadata(*startTime/1000000000, *ip, *country, username)
	case message.TypeAuthPasswordSuccessful:
		u := msg.Payload.(message.PayloadAuthPassword).Username
		username = &u
		storage.SetMetadata(*startTime/1000000000, *ip, *country, username)
	case message.TypeAuthPubKeySuccessful:
		payload := msg.Payload.(message.PayloadAuthPubKey)
		username = &payload.Username
		storage.SetMetadata(*startTime/1000000000, *ip, *country, username)
	case message.TypeHandshakeSuccessful:
		payload := msg.Payload.(message.PayloadHandshakeSuccessful)
		username = &payload.Username
		storage.SetMetadata(*startTime/1000000000, *ip, *country, username)
	}
}
