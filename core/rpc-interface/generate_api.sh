#!/usr/bin/env bash

if [ $(which swagger) == "swagger not found" ];then
    echo "swagger not found, please install swagger or set go bin to path"
    exit 2
fi

if [ $(which java) == "java not found" ];then
    echo "java not found, please install java or set go bin to path"
    exit 2
fi

monitor=./../../../Dipperin-monitor

#echo ${monitor}
if [ ! -f ${monitor}/util/api-util/swagger2markup-cli-1.3.3.jar ];then
    echo 'swagger2markup-cli-1.3.3.jar not found, need copy jar file to path'
    echo 'https://jcenter.bintray.com/io/github/swagger2markup/swagger2markup-cli/1.3.3/'
    echo "cp swagger2markup-cli-1.3.3.jar to ${monitor}/util/api-util/"
    exit 2
fi

generate_api_md () {
    targetDir=$(cd $(dirname $0); pwd)
    echo ${targetDir}
    cd ${targetDir}

    if [ ! -d ${targetDir}/api ];then
        mkdir -p ${targetDir}/api
    fi

    echo "generate swagger json .... "
    swagger generate spec -o ./api/api-doc.json --scan-models

    if [ ! -f ${targetDir}/api/api-doc.json ];then
        echo 'api-doc.json not found, need generate swagger json'
        exit 2
    fi

    echo "convert swagger json to markdown .... "
    java -jar ${root}/util/api-util/swagger2markup-cli-1.3.3.jar convert -i api/api-doc.json -f api/api -c ${root}/util/api-util/config.properties

    cd ${root}
}

generate_api_md