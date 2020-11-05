package s3

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/containerssh/log"

	"github.com/containerssh/auditlog/storage"
)

// NewStorage Creates a storage driver for an S3-compatible object storage.
func NewStorage(cfg Config, logger log.Logger) (storage.ReadWriteStorage, error) {
	httpClient, err := getHTTPClient(cfg)
	if err != nil {
		return nil, err
	}

	awsConfig, partSize, parallelUploads, err2 := getAWSConfig(cfg, logger, httpClient)
	if err2 != nil {
		return nil, err2
	}

	sess := session.Must(session.NewSession(awsConfig))

	queue := newUploadQueue(
		cfg.Local,
		partSize,
		parallelUploads,
		cfg.Bucket,
		cfg.ACL,
		cfg.Metadata.Username,
		cfg.Metadata.IP,
		sess,
		logger,
	)

	if _, err := os.Stat(cfg.Local); err != nil {
		return nil, fmt.Errorf("invalid local audit directory %s (%w)", cfg.Local, err)
	}

	if err := filepath.Walk(cfg.Local, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Size() > 0 && !strings.Contains(info.Name(), ".") {
			if err := queue.recover(info.Name()); err != nil {
				return fmt.Errorf("failed to enqueue old audit log file %s (%w)", info.Name(), err)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return queue, nil
}

func getAWSConfig(
	cfg Config, logger log.Logger, httpClient *http.Client,
) (
	*aws.Config, uint, uint, error,
) {
	var endpoint *string
	if cfg.Endpoint != "" {
		endpoint = &cfg.Endpoint
	}

	if cfg.Bucket == "" {
		return nil, 0, 0, fmt.Errorf("no bucket name specified")
	}
	if cfg.Region == "" {
		return nil, 0, 0, fmt.Errorf("no region name specified")
	}

	awsConfig := &aws.Config{
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     cfg.AccessKey,
				SecretAccessKey: cfg.SecretKey,

				SessionToken: "",
				ProviderName: "",
			},
		}),
		Endpoint:         endpoint,
		Region:           &cfg.Region,
		HTTPClient:       httpClient,
		Logger:           logger,
		S3ForcePathStyle: aws.Bool(cfg.PathStyleAccess),
	}

	partSize := uint(5242880)
	if cfg.UploadPartSize > 5242880 {
		partSize = cfg.UploadPartSize
	}
	parallelUploads := uint(20)
	if cfg.ParallelUploads > 1 {
		parallelUploads = cfg.ParallelUploads
	}
	return awsConfig, partSize, parallelUploads, nil
}

func getHTTPClient(cfg Config) (*http.Client, error) {
	httpClient := http.DefaultClient
	if cfg.CaCert != "" {
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
		if ok := rootCAs.AppendCertsFromPEM([]byte(cfg.CaCert)); !ok {
			return nil, fmt.Errorf("failed to add certificate from config file")
		}
		tlsConfig := &tls.Config{
			RootCAs: rootCAs,
		}
		httpTransport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		httpClient = &http.Client{
			Transport:     httpTransport,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		}
	}
	return httpClient, nil
}
