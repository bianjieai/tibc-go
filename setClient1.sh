chainId="testA" # 本链ID
homePath="${chainId}/node0/simd" # 链执行路径，默认为此路径
address="" # node0 账户地址
genesisChainName="\"native_chain_name\": \"testCreateClientA\"," # 本链轻客户端的名字
voteTime="30s" # 投票时间
genesisVotePeriod="\"voting_period\": \"${voteTime}\"" # 投票时间

rm -rf ${chainId}

make install     # install simd
simd testnet --v 1 -o ${chainId} --chain-id ${chainId} # creat testnet
address=$(simd keys show node0 -a --home ${homePath}) # get keys address
gsed -i "11c minimum-gas-prices =\"0stake\"" ${homePath}/config/app.toml  # change gas to 0
gsed -i "191c ${genesisVotePeriod}" ${homePath}/config/genesis.json # change votePeriod 
gsed -i "249c ${genesisChainName}"  ${homePath}/config/genesis.json  # change chainName 

# 更改端口号
# config
gsed -i "15c proxy_app = \"tcp://127.0.0.1:36658\"" ${homePath}/config/config.toml # change 26658 => 36658
gsed -i "91c laddr = \"tcp://0.0.0.0:36657\"" ${homePath}/config/config.toml # change 26657 => 36657
gsed -i "167c pprof_laddr = \"localhost:6061\"" ${homePath}/config/config.toml # change 6060 => 6061
gsed -i "175c laddr = \"tcp://0.0.0.0:36656\"" ${homePath}/config/config.toml # change 26656 => 36656
gsed -i "392c prometheus_listen_addr = \":26661\"" ${homePath}/config/config.toml # change 26660 => 26661
# app
gsed -i "111c address = \"tcp://0.0.0.0:1318\"" ${homePath}/config/app.toml # change 1317 => 1318
gsed -i "138c address = \":8180\"" ${homePath}/config/app.toml # change 8080 => 8180
gsed -i "162c address = \"0.0.0.0:9190\"" ${homePath}/config/app.toml # change 9090 => 9190
gsed -i "175c address = \"0.0.0.0:9191\"" ${homePath}/config/app.toml # change 9091 => 9191


simd start --home ${homePath}  >  ${homePath}/simd.log 2>&1 & # log in simd.log 
simd keys list --home ${homePath}

exit 0


