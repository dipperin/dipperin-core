#!/usr/bin/env bash

#####################################################################################
#                                                                                   #
#                                                                                   #
#                                                                                   #
#                               Dipperin Shell Script                             #
#                                                                                   #
#                                                                                   #
#                                                                                   #
#####################################################################################

root=`pwd`

GoPath=$HOME/go/src/

cover () {
target_dir=${root}/test-summary

# mkdir
if [[ ! -d ${target_dir} ]];then
    echo "mkdir test-summary"
    mkdir -p ${target_dir}
else
    rm -rf test-summary
    mkdir -p ${target_dir}
fi

# run go test
echo "summary test coverage for dipperin core & cmd & common"
run_c

# summary
cd ${root}/test-summary

go tool cover -func=test_cover.out -o test_summary.text

LASTLINE=$(tail -1 test_summary.text)

echo "core & cmd & common test coverage : " ${LASTLINE}

if [[ $1 == true ]];then
    go tool cover -html=test_cover.out
fi

}

run_c () {
#go test -failfast -v -coverprofile=./test-summary/test_cover.out -p 1 ./core/... ./common/... ./cmd/...
    go test -failfast -coverprofile=./test-summary/test_cover.out -p 1 ./core/... ./cmd/... ./common/...
#go test -failfast -coverprofile=./test-summary/test_cover.out -p 1 ./cmd/...

FAIL=$?

if [[ ${FAIL} != 0 ]];then
    echo "go test failed, please fix test"
    exit ${FAIL}
fi
}

run_test () {

#    is_jenkins="yes" go test -p 1 -cover ./... -coverprofile=size_coverage.out
    is_jenkins="yes" go test -p 1 -cover ./core/... ./common/... ./cmd/...

    return $?

}

jks_test () {

    if [ ! -d ${GoPath} ];then
        echo "mkdir go path:$GoPath"
        mkdir -p ${GoPath}
    fi

    # 获取当前项目的信息
    curProName=`basename $PWD`
    # 获取该项目应该在的位置（go path下）
    gpCurPro="$GoPath$curProName"

    if [ ! -d ${gpCurPro} ];then
        echo "ln cur pro to go path"
        curProPath=`pwd`
        ln -s ${curProPath} ${gpCurPro}
    fi

    cd ${gpCurPro}

    run_test

    return $?

}

build_ci() {
    echo 'build dipperin'

    monitor_path="${GoPath}/chainstack-monitor/core"
    #rm ~/go/bin/dipperin;
    cd ${monitor_path}/cmd/dipperin; go install

    #rm ~/go/bin/dipperincli;
    echo 'build dipperincli'
    cd ${monitor_path}/cmd/dipperincli; go install

    #rm ~/go/bin/bootnode;
    echo 'build bootnode'
    cd ${root}/cmd/bootnode; go install

    #rm ~/go/bin/miner;
    echo 'build miner'
    cd ${root}/cmd/miner; go install

    cd ${root}

    return $?
}


build_install() {
    echo 'build dipperin'

    cd ${root}/cmd/dipperin; go install

    echo 'build dipperincli'
    cd ${root}/cmd/dipperincli; go install

    echo 'build bootnode'
    cd ${root}/cmd/bootnode; go install

    echo 'build miner'
    cd ${root}/cmd/miner; go install

#    echo 'build chain_checker'
#    cd ${root}/cmd/chain_checker; go install

    cd ${root}

    return $?
}

cross_compile() {
    echo 'cross compile dipperin start'
    docker pull karalabe/xgo-latest
    go get github.com/karalabe/xgo
    cd ./cmd/dipperin/
    xgo -go 1.11.1 --targets=linux/amd64,windows/amd64,darwin/amd64 .
    echo 'cross compile dipperin end'
    ls
}

travis_test() {
    echo 'travis test dipperin'
    cache=$(go env | grep "GOCACHE")
    removePath=${cache#*=}
    finalPath=${removePath:1:-1}

    if [ "$removePath" != "" ];then
        echo "remove the GOCACHE"
        echo $finalPath
        #rm -rf $finalPath
    fi

    go test ./...
}

update_vendor () {


# enable go1.11 mod in go path
export GO111MODULE=on

if [[ -n "$1" ]];then

    ############################################
    #
    #
    #            need use http proxy
    #
    #
    ############################################

    # go mod tidy download pkg which we are using
    export http_proxy=http://127.0.0.1:$1;export https_proxy=http://127.0.0.1:$1;go mod tidy
    go mod vendor

    exit $?
fi

go mod tidy

go mod vendor

return $?
}

main () {
 if [[ ! $1 ]];then
        run_test
        exit $?
 fi

 case $1 in

        [Jj][Kk][Ss] )
            {
                jks_test
        };;

        ci )
            {
                build_ci
        };;

        install )
            {
                build_install
        };;

        tidy )
            {
                update_vendor $2
        };;

        cover )
            {
                cover $2
        };;

        compile )
            {
                cross_compile
        };;

        travisTest )
            {
                travis_test
        };;

    esac

    exit $?

}

echo "

#####################################################################################
#                                                                                   #
#                                                                                   #
#                                                                                   #
#                               Dipperin Shell Script                             #
#                                                                                   #
#                                                                                   #
#                                                                                   #
#####################################################################################

"
# setup git hooks
find .git/hooks -type l -exec rm {} \;
find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

main $1 $2