module github.com/edgexfoundry/device-protocol

require (
	github.com/edgexfoundry/device-sdk-go v1.0.0
	github.com/edgexfoundry/go-mod-core-contracts v0.1.0
	github.com/edgexfoundry/go-mod-registry v0.1.0
	github.com/google/uuid v1.1.0
	github.com/gorilla/mux v1.6.2
	github.com/pelletier/go-toml v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	github.com/ugorji/go v1.1.4
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20181029021203-45a5f77698d3
	golang.org/x/net => github.com/golang/net v0.0.0-20181023162649-9b4f9f5ad519
	golang.org/x/sync => github.com/golang/sync v0.0.0-20181221193216-37e7f081c4d4
	golang.org/x/sys => github.com/golang/sys v0.0.0-20181026203630-95b1ffbd15a5
	golang.org/x/tools => github.com/golang/tools v0.0.0-20181112210238-4b1f3b6b1646
)
