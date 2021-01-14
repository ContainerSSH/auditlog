module github.com/containerssh/auditlog

go 1.14

require (
	github.com/Microsoft/go-winio v0.4.16 // indirect
	github.com/aws/aws-sdk-go v1.36.26
	github.com/containerssh/geoip v0.9.4
	github.com/containerssh/log v0.9.9
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0 // indirect
	github.com/fxamacker/cbor v1.5.1
	github.com/mitchellh/mapstructure v1.4.1
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/stretchr/testify v1.6.1
)

replace (
	github.com/davecgh/go-spew v1.1.0 => github.com/davecgh/go-spew v1.1.1
	github.com/stretchr/testify v1.2.2 => github.com/stretchr/testify v1.6.1
	github.com/stretchr/testify v1.4.0 => github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 => golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20190404232315-eb5bcb51f2a3 => golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
	golang.org/x/sys v0.0.0-20180905080454-ebe1bf3edb33 => golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/sys v0.0.0-20190412213103-97732733099d => golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/sys v0.0.0-20190916202348-b4ddaad3f8a3 => golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/sys v0.0.0-20191224085550-c709ea063b76 => golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f
	golang.org/x/text v0.3.0 => golang.org/x/text v0.3.3
	gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 => gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f
	gopkg.in/yaml.v2 v2.2.2 => gopkg.in/yaml.v2 v2.3.0
	gopkg.in/yaml.v2 v2.2.8 => gopkg.in/yaml.v2 v2.3.0

)
