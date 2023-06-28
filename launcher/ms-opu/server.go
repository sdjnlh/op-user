package main

import (
	_ "github.com/lib/pq"
	"github.com/sdjnlh/communal/app"
	"github.com/sdjnlh/op-user/opu"
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
