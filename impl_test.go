package auditlog_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"testing"
	"time"

	"github.com/containerssh/log"
	"github.com/containerssh/log/formatter/ljson"
	logPipeline "github.com/containerssh/log/pipeline"
	"github.com/stretchr/testify/assert"

	"github.com/containerssh/auditlog"
	"github.com/containerssh/auditlog/codec/binary"
	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage"
	"github.com/containerssh/auditlog/storage/file"
)

func newTestCase(t *testing.T) (*testCase, error) {
	var err error
	dir, err := ioutil.TempDir(os.TempDir(), "containerssh-auditlog-test")
	if err != nil {
		t.Skipf("failed to create temporary directory (%v)", err)
		return nil, err
	}
	tc := &testCase{
		dir: dir,
		t:   t,
		config: auditlog.Config{
			Format:  "binary",
			Storage: "file",
			File: file.Config{
				Directory: dir,
			},
			Intercept: auditlog.InterceptConfig{
				Stdin:     true,
				Stdout:    true,
				Stderr:    true,
				Passwords: true,
			},
		},
	}
	if err := tc.setUpLogger(); err != nil {
		assert.Fail(t, "failed to set up test case", err)
		return nil, err
	}
	return tc, nil
}

type testCase struct {
	dir         string
	t           *testing.T
	auditLogger auditlog.Logger
	logger      *logPipeline.LoggerPipeline
	config      auditlog.Config
}

func (c *testCase) setUpLogger() error {
	c.logger = logPipeline.NewLoggerPipeline(log.LevelDebug, ljson.NewLJsonLogFormatter(), os.Stdout)
	auditLogger, err := auditlog.New(c.config, c.logger)
	if err != nil {
		assert.Fail(c.t, "failed to create audit logger", err)
		return err
	}
	c.auditLogger = auditLogger
	return nil
}

func (c *testCase) listAuditLogs() ([]storage.Entry, error) {
	fileStorage, err := file.NewStorage(c.config.File, nil)
	if err != nil {
		return nil, err
	}
	var logs []storage.Entry
	logsChannel, errors := fileStorage.List()
	for {
		finished := false
		select {
		case entry, ok := <-logsChannel:
			if !ok {
				finished = true
				break
			}
			logs = append(logs, entry)
		case err, ok := <-errors:
			if !ok {
				finished = true
				break
			}
			if err != nil {
				return nil, err
			}
		}
		if finished {
			break
		}
	}
	return logs, nil
}

func (c *testCase) getRecentAuditLogMessages(t *testing.T) ([]message.Message, error) {
	auditLogs, err := c.listAuditLogs()
	if err != nil {
		assert.Fail(t, "failed to list audit logs", err)
	}

	if !assert.Equal(t, 1, len(auditLogs)) {
		return nil, fmt.Errorf("invalid number of audit logs")
	}

	messages, err := c.getAuditLog(auditLogs[0].Name)
	if err != nil {
		assert.Fail(t, "failed to fetch audit log", err)
		return nil, fmt.Errorf("failed to fetch audit log")
	}
	return messages, nil
}

func (c *testCase) tearDown() {
	err := os.RemoveAll(c.dir)
	if err != nil {
		// Give the audit logger time to close the file on Windows
		time.Sleep(2 * time.Second)
		err2 := os.RemoveAll(c.dir)
		if err2 != nil {
			assert.Fail(c.t, "failed to remove temporary directory after test", err)
		}
	}
}

func (c *testCase) getAuditLog(name string) ([]message.Message, error) {
	fileStorage, err := file.NewStorage(c.config.File, nil)
	if err != nil {
		return nil, err
	}
	reader, err := fileStorage.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	decoder := binary.NewDecoder()
	messageChannel, errors := decoder.Decode(reader)
	var result []message.Message
	for {
		finished := false
		select {
		case msg, ok := <-messageChannel:
			if !ok {
				finished = true
				break
			}
			result = append(result, msg)
		case err, ok := <-errors:
			if !ok {
				finished = true
				break
			}
			return nil, err
		}
		if finished {
			break
		}
	}
	return result, nil
}

func TestConnect(t *testing.T) {
	testCase, err := newTestCase(t)
	if err != nil {
		//Already handled
		return
	}
	defer testCase.tearDown()
	auditLogger := testCase.auditLogger

	connection, err := auditLogger.OnConnect(
		[]byte("asdf"),
		net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 2222,
			Zone: "",
		},
	)
	if err != nil {
		assert.Fail(t, "failed to send connect message to logger", err)
		return
	}
	connection.OnDisconnect()

	messages, err := testCase.getRecentAuditLogMessages(t)
	if err != nil {
		return
	}

	assert.Equal(t, 2, len(messages))

	assert.Equal(t, message.TypeConnect, messages[0].MessageType)
	assert.Equal(t, []byte("asdf"), []byte(messages[0].ConnectionID))
	payload1 := messages[0].Payload.(message.PayloadConnect)
	assert.Equal(t, "127.0.0.1", payload1.RemoteAddr)

	assert.Equal(t, message.TypeDisconnect, messages[1].MessageType)
	assert.Equal(t, []byte("asdf"), []byte(messages[1].ConnectionID))
	assert.Equal(t, nil, messages[1].Payload)
}

func TestAuth(t *testing.T) {
	testCase, err := newTestCase(t)
	if err != nil {
		//Already handled
		return
	}
	defer testCase.tearDown()

	connection, err := testCase.auditLogger.OnConnect(
		[]byte("asdf"),
		net.TCPAddr{
			IP:   net.ParseIP("127.0.0.1"),
			Port: 2222,
			Zone: "",
		},
	)
	if err != nil {
		assert.Fail(t, "failed to send connect message to logger", err)
		return
	}
	connection.OnAuthPassword("foo", []byte("bar"))
	connection.OnAuthPasswordBackendError("foo", []byte("bar"), "no particular reason")
	connection.OnAuthPassword("foo", []byte("bar"))
	connection.OnAuthPasswordFailed("foo", []byte("bar"))
	connection.OnAuthPassword("foo", []byte("baz"))
	connection.OnAuthPasswordSuccess("foo", []byte("baz"))
	connection.OnAuthPubKey("foo", []byte("bar"))
	connection.OnAuthPubKeyBackendError("foo", []byte("bar"), "no particular reason")
	connection.OnAuthPubKey("foo", []byte("bar"))
	connection.OnAuthPubKeyFailed("foo", []byte("bar"))
	connection.OnAuthPubKey("foo", []byte("baz"))
	connection.OnAuthPubKeySuccess("foo", []byte("baz"))
	connection.OnDisconnect()

	messages, err := testCase.getRecentAuditLogMessages(t)
	if err != nil {
		return
	}

	assert.Equal(t, 14, len(messages))

	for _, msg := range messages {
		assert.Equal(t, []byte("asdf"), []byte(msg.ConnectionID))
	}
	assert.Equal(t, message.TypeAuthPassword, messages[1].MessageType)
	assert.Equal(t, message.TypeAuthPasswordBackendError, messages[2].MessageType)
	assert.Equal(t, message.TypeAuthPassword, messages[3].MessageType)
	assert.Equal(t, message.TypeAuthPasswordFailed, messages[4].MessageType)
	assert.Equal(t, message.TypeAuthPassword, messages[5].MessageType)
	assert.Equal(t, message.TypeAuthPasswordSuccessful, messages[6].MessageType)
	assert.Equal(t, message.TypeAuthPubKey, messages[7].MessageType)
	assert.Equal(t, message.TypeAuthPubKeyBackendError, messages[8].MessageType)
	assert.Equal(t, message.TypeAuthPubKey, messages[9].MessageType)
	assert.Equal(t, message.TypeAuthPubKeyFailed, messages[10].MessageType)
	assert.Equal(t, message.TypeAuthPubKey, messages[11].MessageType)
	assert.Equal(t, message.TypeAuthPubKeySuccessful, messages[12].MessageType)
}
