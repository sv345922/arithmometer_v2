package app

import (
	"context"
	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/useCases/getResult"
	"github.com/sv345922/arithmometer_v2/internal/useCases/getTask"
	"github.com/sv345922/arithmometer_v2/internal/useCases/giveAnswer"
	"github.com/sv345922/arithmometer_v2/internal/useCases/newExpression"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
	"log"
	"net/http"
)

func RunServer(ctx context.Context, ws *wSpace.WorkingSpace) error {
	mux := http.NewServeMux()

	// Дать ответ клиенту о результатах вычисления выражений
	mux.HandleFunc("/getresult", getResult.GetResult(ws)) // TODO context?

	// Получение нового выражения от клиента
	mux.HandleFunc("/newexpression", newExpression.NewExpression(ctx, ws))

	// Дать задачу вычислителю
	mux.HandleFunc("/gettask", getTask.GetTask(ctx, ws)) // TODO context

	// Получить ответ от вычислителя
	mux.HandleFunc("/giveanswer", giveAnswer.GiveAnswer(ctx, ws)) // TODO context

	log.Println("Server is working")
	defer log.Println("Server stopped")
	err := http.ListenAndServe("localhost:"+configs.Port, mux)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}
