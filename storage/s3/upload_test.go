package s3_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	goLog "log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/containerssh/log"
	"github.com/containerssh/log/formatter/ljson"
	"github.com/containerssh/log/pipeline"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"

	auditLogStorage "github.com/containerssh/auditlog/storage"
	"github.com/containerssh/auditlog/storage/s3"
)

type minio struct {
	containerID string
	dir         string
	storage     auditLogStorage.ReadWriteStorage
}

func (m *minio) getClient() (*client.Client, error) {
	dockerURL := "/var/run/docker.sock"
	_, err := os.Stat(dockerURL)
	if err != nil {
		dockerURL = "tcp://127.0.0.1:2375"
	}

	return client.NewClient(dockerURL, "", nil, make(map[string]string))
}

func (m *minio) Start(
	t *testing.T,
	accessKey string,
	secretKey string,
	region string,
	bucket string,
	endpoint string,
) (auditLogStorage.ReadWriteStorage, error) {
	if m.containerID == "" {
		storage, err2 := m.startMinio(t, accessKey, secretKey)
		if err2 != nil {
			return storage, err2
		}
	}

	var err error
	if err := setupS3(t, accessKey, secretKey, endpoint, region, bucket); err != nil {
		return nil, err
	}

	m.dir, err = ioutil.TempDir(os.TempDir(), "containerssh-s3-upload-test")
	if err != nil {
		t.Skipf("failed to create temporary directory (%v)", err)
		return nil, err
	}

	logger := pipeline.NewLoggerPipeline(log.LevelDebug, ljson.NewLJsonLogFormatter(), os.Stdout)
	m.storage, err = s3.NewStorage(
		s3.Config{
			Local:           m.dir,
			AccessKey:       accessKey,
			SecretKey:       secretKey,
			Bucket:          bucket,
			Region:          region,
			Endpoint:        endpoint,
			PathStyleAccess: true,
			Metadata:        s3.Metadata{},
		},
		logger,
	)
	if err != nil {
		assert.Fail(t, "failed to create storage (%v)", err)
		return nil, err
	}

	return m.storage, nil
}

func (m *minio) startMinio(t *testing.T, accessKey string, secretKey string) (auditLogStorage.ReadWriteStorage, error) {
	ctx := context.Background()

	cli, err := m.getClient()
	if err != nil {
		t.Skipf("failed to create Docker client (%v)", err)
		return nil, err
	}

	reader, err := cli.ImagePull(ctx, "docker.io/minio/minio", types.ImagePullOptions{})
	if err != nil {
		t.Skipf("failed to pull Minio image (%v)", err)
		return nil, err
	}
	if _, err := io.Copy(os.Stdout, reader); err != nil {
		t.Skipf("failed to stream logs from Minio image pull (%v)", err)
		return nil, err
	}

	env := []string{
		fmt.Sprintf("MINIO_ACCESS_KEY=%s", accessKey),
		fmt.Sprintf("MINIO_SECRET_KEY=%s", secretKey),
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: "minio/minio",
			Cmd:   []string{"server", "/data"},
			Env:   env,
		},
		&container.HostConfig{
			PortBindings: map[nat.Port][]nat.PortBinding{
				"9000/tcp": {
					{
						HostIP:   "127.0.0.1",
						HostPort: "9000",
					},
				},
			},
		},
		nil,
		"",
	)
	if err != nil {
		t.Skipf("failed to create Minio container (%v)", err)
		return nil, err
	}

	m.containerID = resp.ID

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Skipf("failed to start minio container (%v)", err)
		return nil, err
	}
	return nil, nil
}

func (m *minio) Stop(t *testing.T) {
	if m.containerID != "" {
		ctx := context.Background()

		cli, err := m.getClient()
		if err != nil {
			t.Skipf("failed to create Docker client (%v)", err)
		}

		if err = cli.ContainerRemove(ctx, m.containerID, types.ContainerRemoveOptions{
			RemoveVolumes: false,
			RemoveLinks:   false,
			Force:         true,
		}); err != nil {
			goLog.Println("failed to remove Minio container")
		}
		m.containerID = ""
	}

	if m.dir != "" {
		m.dir = ""
		if err := os.RemoveAll(m.dir); err != nil {
			goLog.Printf("failed to remove temporary directory (%v)", err)
		}
	}
}

func setupS3(t *testing.T, accessKey string, secretKey string, endpoint string, region string, bucket string) error {
	awsConfig := &aws.Config{
		CredentialsChainVerboseErrors: nil,
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     accessKey,
				SecretAccessKey: secretKey,
			},
		}),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		S3ForcePathStyle: aws.Bool(true),
	}
	sess := session.Must(session.NewSession(awsConfig))
	s3Connection := awsS3.New(sess)
	if _, err := s3Connection.CreateBucket(&awsS3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		t.Skipf("failed to create bucket (%v)", err)
		return err
	}
	return nil
}

func getS3Objects(t *testing.T, storage auditLogStorage.ReadWriteStorage) []auditLogStorage.Entry {
	var objects []auditLogStorage.Entry
	objectChan, errChan := storage.List()
	for {
		finished := false
		select {
		case object, ok := <-objectChan:
			if !ok {
				finished = true
				break
			}
			objects = append(objects, object)
		case err, ok := <-errChan:
			if !ok {
				finished = true
				break
			}
			assert.Fail(t, "error while fetching objects from object storage", err)
		}
		if finished {
			break
		}
	}
	return objects
}

func waitForS3Objects(t *testing.T, storage auditLogStorage.ReadWriteStorage, count int) []auditLogStorage.Entry {
	tries := 0
	var objects []auditLogStorage.Entry
	for {
		if tries > 10 {
			break
		}
		objects = getS3Objects(t, storage)
		if len(objects) > count-1 {
			break
		} else {
			tries++
			time.Sleep(10 * time.Second)
		}
	}
	return objects
}

func TestSmallUpload(t *testing.T) {
	accessKey := "asdfasdfasdf"
	secretKey := "asdfasdfasdf"

	region := "us-east-1"
	bucket := "auditlog"
	endpoint := "http://127.0.0.1:9000"

	m := &minio{}

	storage, err := m.Start(t, accessKey, secretKey, region, bucket, endpoint)
	if err != nil {
		return
	}
	defer func() {
		m.Stop(t)
	}()

	writer, err := storage.OpenWriter("test")
	if err != nil {
		assert.Fail(t, "failed to open storage writer (%v)", err)
		return
	}
	var data = []byte("Hello world!")
	if _, err := writer.Write(data); err != nil {
		assert.Fail(t, "failed to write to storage writer (%v)", err)
		return
	}
	if err := writer.Close(); err != nil {
		assert.Fail(t, "failed to close storage writer (%v)", err)
		return
	}

	objects := waitForS3Objects(t, storage, 1)
	assert.Equal(t, 1, len(objects))

	r, err := storage.OpenReader(objects[0].Name)
	if err != nil {
		assert.Fail(t, "failed to open reader for recently stored object", err)
		return
	}
	d, err := ioutil.ReadAll(r)
	if err != nil {
		assert.Fail(t, "failed to open read from S3", err)
		return
	}

	assert.Equal(t, data, d)
}
