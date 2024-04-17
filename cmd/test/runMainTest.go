package main

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/dataBase"
	"github.com/sv345922/arithmometer_v2/internal/grps/server"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
)

func main() {
	fmt.Println("main Test running...")
	ctx := context.Background()
	dbase, err := dataBase.CreateDb(ctx, "testDB")
	if err != nil {
		panic(err)
	}
	defer dbase.Close()
	ws, err := wSpace.LoadDB(ctx, dbase)
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	//запуск сервера grpc
	server.StartGRPCServer(ctx, ws)

	//err = app.RunServer(ctx, ws)
	//log.Println(err)

}
