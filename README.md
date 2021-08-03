# tibc-go
Golang Implementation of Terse IBC

## build
```bash
make build
```

## local testnet
```bash
./build/simd testnet --v 1 --chain-id test --keyring-backend file

./build/simd start --home mytestnet/node0/simd
```