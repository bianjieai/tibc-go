package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/bianjieai/tibc-go/simapp"
	"github.com/bianjieai/tibc-go/simapp/simd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "simd", simapp.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
