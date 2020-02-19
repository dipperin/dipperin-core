#!/bin/bash
#add copy right head to the go file
#set -v

skipStr="// Copyright"
skipStr2="//ISC Licens"
DipperinStr="Keychain"

if [ ! $1 ]
then
	echo "please assign file dir"
	exit	 
fi

if [ ! $2 ]
then
	echo "please assign copy right file"
	exit
fi

funcAddCopyRight(){
	echo "handle file: $1"
	value=$(sed -n '1p' $1)
	cmpStr=${value:0:12}	
	
	keyChainStr=${value:22:8}
	
	if [ "$cmpStr" != "$skipStr" ] && [ "$cmpStr" != "$skipStr2" ]
	then
		cat $2 $1 >tmpFile && mv tmpFile $1
	else
		if [ "$keyChainStr" == "$DipperinStr" ]
		then
			sed -i '1,15d' $1
			cat $2 $1 >tmpFile && mv tmpFile $1
		fi	
	fi
}

for file in $(find $1 -name "*.go" ! -path "*/third-party/*" ! -path "*/vendor/*")
do
	funcAddCopyRight $file $2
done

echo "run end"
