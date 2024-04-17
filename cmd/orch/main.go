package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/sv345922/arithmometer_v2/internal/app"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/dataBase"
)

// создать задачу (выражение)
// зафиксировать тайминги операторов
// сохранить выражение в БД
// сделать список задач для вычисления
// сохранить список задач
// отдать задачу вычислителю
// получить ответ от вычислителя
// обновить список задач
// повторить до завершения всех задач
// вернуть ответ клиенту при запросе

func main() {
	// создать пустую базу
	ctx, stopCtx := context.WithCancel(context.Background())
	var dbase *sql.DB
	var err error
	if len(os.Args) > 1 {
		if os.Args[1] == "new" {
			dbase, err = dataBase.CreateEmptyDb(ctx)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		dbase, err = dataBase.CreateDb(ctx, configs.DBPath)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer dbase.Close()

	ws, err := app.RunTasker(ctx, dbase)
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	app, _ := app.New(ctx, ws)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		app.GracefulStop(ctx, sig, stopCtx)
	}()
	err = app.Run()
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	<-ctx.Done()

}
