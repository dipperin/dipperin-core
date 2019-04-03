#!/usr/bin/env bash

##remote_init
#scp -P 3222 ~/go/bin/remote_init dgcs@dgly:~/go/bin

##dipperin
go install
scp -P 3222 ~/go/bin/Dipperin dgcs@dgly:~/go/bin

##bootnode
#cd ../bootnode
#go install
#scp -P 3222 ~/go/bin/bootnode dgcs@dgly:~/go/bin


#scp -P 3222 ~/go/bin/cs_ci_ex dgcs@dgly:~/go/bin
