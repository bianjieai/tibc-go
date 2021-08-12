module github.com/bianjieai/tibc-go

go 1.16

require (
	github.com/armon/go-metrics v0.3.9
	github.com/confio/ics23/go v0.6.6
	github.com/cosmos/cosmos-sdk v0.42.9
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/irisnet/irismod v1.4.1-0.20210810111846-cc9a8f07dc53
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.4.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.11
	github.com/tendermint/tm-db v0.6.4
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.39.0
	google.golang.org/protobuf v1.27.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
