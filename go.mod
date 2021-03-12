module github.com/containerssh/auditlog

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/aws/aws-sdk-go v1.37.29
	github.com/containerd/containerd v1.4.3 // indirect
	github.com/containerssh/geoip v0.9.4
	github.com/containerssh/log v0.9.13
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v20.10.5+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fxamacker/cbor v1.5.1
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/mitchellh/mapstructure v1.4.1
	github.com/moby/term v0.0.0-20201216013528-df9cb8a40635 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/stretchr/testify v1.7.0
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324 // indirect
	google.golang.org/grpc v1.35.0 // indirect
	gotest.tools/v3 v3.0.3 // indirect
)

replace (
	// Fixes CVE-2020-9283
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550 => golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975 => golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 => golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	// Fixes CVE-2020-14040
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.1 => golang.org/x/text v0.3.3
	golang.org/x/text v0.3.2 => golang.org/x/text v0.3.3
)
