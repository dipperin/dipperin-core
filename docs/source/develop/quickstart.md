# Quick start
Welcome to use Dipperin. Follow this guide you can run a Dipperin node on Dipperin Testnet.

## Prerequisites
Dipperin uses the Go Programming Language for many of its components.

- [Go](https://golang.org/dl/) version 1.11.x is required.
- C compiler

## Building the source

### Mac & Linux

Get source code to your go path

```shell
$ mkdir -p ~/go/src/github.com/dipperin
$ cd ~/go/src/github.com/dipperin
$ git clone git@github.com:dipperin/dipperin-core.git
```

Run tests

```shell
$ cd ~/go/src/github.com/dipperin/dipperin-core
$ make test
# or
$ ./cs.sh
```

Build dipperin to `~/go/bin`

```shell
$ cd ~/go/src/github.com/dipperin/dipperin-core
$ make build
```

Build following softwares to `~/go/bin`

1. dipperin
2. dipperincli
3. bootnode
4. miner
5. chain_checker

```
$ ./cs.sh install
```

### Windows

The Chocolatey package manager provides an easy way to get the required build tools installed. If you don't have chocolatey yet, follow the instructions on [https://chocolatey.org](https://chocolatey.org) to install it first.

Then open an Administrator command prompt and install the build tools we need:

```
C:\Windows\system32> choco install git
C:\Windows\system32> choco install golang
C:\Windows\system32> choco install mingw
```

Use git shell run commands below, and copy source code to your go path

You can't run the tests if you don't put source code in your home folder.

(`$HOME` means home folder, example `C:\Users\qydev`)

```
$ mkdir -p $HOME\go\src\github.com\dipperin
$ cd $HOME\go\src\github.com\dipperin
$ git clone git@github.com:dipperin/dipperin-core.git
```

Restart cmd and run tests

```
$ cd $HOME\go\src\github.com\dipperin\dipperin-core
$ go test -p 1 ./...
```

Build dipperin to User

```
$ cd $HOME\go\src\github.com\dipperin\dipperin-core\cmd\dipperin
$ go install
```

Build dipperincli to User

```
$ cd $HOME\go\src\github.com\dipperin\dipperin-core\cmd\dipperincli
$ go install
```

## Executables

The dipperin-core project comes with several wrappers/executables found in the `cmd` directory.

- dipperin

Our chain CLI client. It is the entry point into the Dipperin network, capable of running as a full node.
It can be used by other processes as a gateway into the Dipperin network via JSON RPC endpoints exposed on top of HTTP,
WebSocket and/or IPC transports.

- dipperincli

Our chain CLI client with console. It has all features of `dipperin`, and provides a easy way to start the node.
You can give commands to node in command line console, like starting mining `miner StartMine` or querying current block `chain CurrentBlock`

- bootnode

Stripped down version of our Dipperin client implementation that only takes part in the network node discovery protocol,
but does not run any of the higher level application protocols.
It can be used as a lightweight bootstrap node to aid in finding peers in private networks.

- miner

Mine block client, It must work with a `mine master` started by `dipperin` or `dipperincli`.
`mine master` dispatch sharding works for every miner registered.
So all miner do different works when mining a block. |

## Running dipperin

### Setting environment variables

- Mac & Linux

```
$ vi ~/.bashrc
```

Add command `export PATH=$PATH:~/go/bin` at bottom of the file,
and `:wq` for save and quit the file.

- Windows

(`$HOME` means home folder, example `C:\Users\qydev`)

```
$ set PATH=%PATH%;$HOME\go\bin
```

Going through all the possible command line flags

```shell
$ dipperin -h
# or
$ dipperincli -h
```

### Full node on the main Dipperin network

- Mac & Linux

```shell
$ boots_env=venus dipperincli
```

- Windows

```shell
$ set boots_env=venus
$ dipperincli
```

This command will:

 * Guide you to setup your personal Dipperin start config, and will write these args to your `$HOME/.dipperin/start_conf.json`, you can change these start args in this file.
 * Start sync Dipperin test-net data from other nodes.

## Using command line

The following command is to start a node, which requires a wallet password.

If no wallet path is specified, the default system path is used: `$Home/.dipperin/`.

```
$ dipperincli -- node_type [type] -- soft_wallet_pwd [password]
```

Example:

Local startup miner:
```
$ dipperincli -- node_type 1 -- soft_wallet_pwd 123
```

Local startup miner(start mining):
```
$ dipperincli -- node_type 1 -- soft_wallet_pwd 123 -- is_start_mine 1
```

Local startup verifier:
```
$ dipperincli -- node_type 2 -- soft_wallet_pwd 123
```

### Error

If dipperincli started in a wrong way,
it may be that the local link data is not synchronized with the link state,
and the local link data needs to be deleted:

- Mac & Linux

```
$ cd ~
$ rm .dipperin -fr
```

- Windows

```
$ rd /s /q $HOME\.dipperin
```

restart command line tool

[See more details for Command Line Tool](../design/commands.html)
