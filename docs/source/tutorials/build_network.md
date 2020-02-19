# Build your first network

## Install and start Dipperin Command Line Tool

If you don't know how to install or start Dipperin Command Line Tool, please take a look at the [Quick Start](../develop/quickstart)

## Init verifiers

After start verifier nodes, you'll see

```shell
t=2019-01-27T10:44:46+0800 lvl=info msg="setup default sign address" addr=0x0000C82ADd56D1E464719D606bB873Ad52779c67F465
```

copy this `address` like `0x0000C82ADd56D1E464719D606bB873Ad52779c67F465` to your genesis verifiers.

```shell
$ dipperin --node_type 2 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/verifier1 --http_port 10001 --ws_port 10002 --p2p_listener 20001
$ dipperin --node_type 2 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/verifier2 --http_port 10003 --ws_port 10004 --p2p_listener 20002
$ dipperin --node_type 2 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/verifier3 --http_port 10005 --ws_port 10006 --p2p_listener 20003
...
```

This is done so that you can generate the default verifier wallet, which you can configure in `genesis.json`.

## Setup genesis state

You need input content below into file `$HOME/softwares/dipperin_deploy/genesis.json`.

```json
{
  "nonce": 11,
  "accounts": {
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af19": 1000,
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af29": 1000,
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af39": 1000
  },
  "timestamp": "1548554091989871000",
  "difficulty": "0x1e566611",
  "verifiers": [
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af19",
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af29",
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af39",
    "0x00005EE98a9d6776F4599f8cD9070843E6D03Ce6af49"
  ]
}
```

In the json, `accounts` is pre-fund some accounts for your private chain. `verifiers` is first round default verifiers for you private chain, this list must have `22` verifiers, you can change this number in `core/chain-config/config.go` at func `defaultChainConfig` -> `VerifierNumber`.

## Start a bootnode

Generate bootnode private key file, and start it.

```shell
$ bootnode --genkey=boot.key
$ bootnode --nodekey=boot.key
```

You'll see the following code:

```

bootnode conn: enode://958784048f7021c99b5ce82bd0078398037226ffd35c166b874fc8ff36d0c4e07e0a2a28eb02b6d993ec8b652f79a9bf79725fcf7ba754bf4c2f670f330b9080@127.0.0.1:30301
```

when bootnode started, copy this conn str to `core/chain-config/config.go` at func `initLocalBoots` -> `KBucketNodes`,
and recompile your `dipperin`, your node will auto connect this bootnode when started.
Or you can write this conn str to your node's `static_boot_nodes.json` file in `datadir`, it's content should like:

```json
[
  "enode://958784048f7021c99b5ce82bd0078398037226ffd35c166b874fc8ff36d0c4e07e0a2a28eb02b6d993ec8b652f79a9bf79725fcf7ba754bf4c2f670f330b9080@127.0.0.1:30301"
]
```

## Start verifiers

You should remove `full_chain_data` in all `datadir` because of your genesis block has changed, and don't remove `CSWallet` in `datadir`.
Then run commands below to started verifiers.

```shell
$ dipperin --node_type 2 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/verifier1 --http_port 10001 --ws_port 10002 --p2p_listener 20001
$ dipperin --node_type 2 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/verifier2 --http_port 10003 --ws_port 10004 --p2p_listener 20002
$ dipperin --node_type 2 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/verifier3 --http_port 10005 --ws_port 10006 --p2p_listener 20003
...
```

## Start miner master(default have a miner)

```shell
$ dipperin --node_type 1 --soft_wallet_pwd 123 --data_dir /home/qydev/dipperin/mine_master1 --http_port 10010 --ws_port 10011 --p2p_listener 20010
...
```

This command will start a `mine master` and start a `miner` in it, you'll see it is mining block and broadcast block to verifiers.

And your private chain block height is growing up.
