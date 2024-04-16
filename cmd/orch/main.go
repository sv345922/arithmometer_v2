package main

import (
	"context"
	"database/sql"
	"github.com/sv345922/arithmometer_v2/internal/app"
	"github.com/sv345922/arithmometer_v2/internal/dataBase"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"os"
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
	ctx := context.Background()
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
		dbase, err = dataBase.CreateDb(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer dbase.Close()
	// сделать список задач для вычисления
	ws, err := RunTasker(ctx, dbase)
	if err != nil {
		log.Fatalf("main: %v", err)
	}
	err = app.RunServer(ctx, ws)
	log.Println(err)
}

// Создает рабочее пространство из сохраненной базы данных
func RunTasker(ctx context.Context, db *sql.DB) (*wSpace.WorkingSpace, error) {
	// Восстанавливаем выражения и задачи из базы данных
	// Загрузка сохраненной БД
	ws, err := wSpace.LoadDB(ctx, db)
	if err != nil {
		// log.Println("ошибка загрузки БД", err)
		return ws, err
	}

	return ws, err
}
