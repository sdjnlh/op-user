workingSpace=$(cd $(dirname $0); pwd)

serviceDir=$(ls launcher|grep ms-)
webDir=$(ls launcher|grep web-)
echo $serviceDir
echo $webDir
sd="$workingSpace/launcher/$serviceDir"
wd="$workingSpace/launcher/$webDir"
#echo $wd
#echo $#

# if [ $# -lt 1 ]; then
#     echo "error.. need args"
#     exit 1
# fi

type="a"

if [ $# -ne 0 ]; then
	type="$1"
	# echo "error.. need args"
fi

# echo $type
# echo $fsd

service(){
	cd $sd
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -tags "consul"
	echo "success build service"
}

web(){
	cd $wd
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -tags "consul jsoniter"
	echo "success build web"
}

if [ $type == "a" ]; then
	# echo "build all"
	service
	web
	echo "done"
elif [ $type == "s" ]; then
	# echo "build service"
	service
	echo "done"
elif [ $type == "w" ]; then
	# echo "build web"
	web
else
	# echo "build all"
	service
	web
	echo "done"
fi
