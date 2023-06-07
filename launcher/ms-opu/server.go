package main

import (
	"code.letsit.cn/go/common/app"
	"code.letsit.cn/go/op-user/opu"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	app.RegisterStarter(opu.Service)
	if err = app.Start(); err != nil {
		panic(err)
	}

	//server := server.NewServer()
	//registry, err := rpc.RegisterConsul(server)

	if err != nil {
		panic(err)
	}

	//service.Register(server)
	//if err = server.Serve("tcp", registry.Client); err != nil {
	//	panic(err)
	//}
}
