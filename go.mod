module github.com/bianjieai/tibc-go

go 1.16

require (
	cosmossdk.io/math v1.0.0-beta.2
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/armon/go-metrics v0.3.11
	github.com/confio/ics23/go v0.7.0
	github.com/cosmos/cosmos-sdk v0.46.0-rc1
	github.com/edsrzf/mmap-go v1.0.0
	github.com/ethereum/go-ethereum v1.10.17
	github.com/gogo/protobuf v1.3.3
	github.com/golang/protobuf v1.5.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/irisnet/irismod v1.5.3-0.20220618121128-2743ff366d16
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/spf13/cast v1.5.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	github.com/tendermint/tendermint v0.35.4
	github.com/tendermint/tm-db v0.6.7
	golang.org/x/crypto v0.0.0-20220411220226-7b82a4e95df4
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
