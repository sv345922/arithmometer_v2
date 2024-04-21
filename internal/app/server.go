package app

import (
	"context"
	"fmt"
	"github.com/sv345922/arithmometer_v2/internal/useCases/authorization"
	"github.com/sv345922/arithmometer_v2/internal/useCases/loginUser"
	"github.com/sv345922/arithmometer_v2/internal/useCases/registration"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sv345922/arithmometer_v2/internal/configs"
	"github.com/sv345922/arithmometer_v2/internal/grps/grpcServer"
	"github.com/sv345922/arithmometer_v2/internal/useCases/getResult"
	"github.com/sv345922/arithmometer_v2/internal/useCases/newExpression"
	"github.com/sv345922/arithmometer_v2/internal/wSpace"
)

type App struct {
	httpServer *http.Server
	gRPCServer *grpcServer.Server
	ctx        context.Context
	ws         *wSpace.WorkingSpace
}

// Создает и настраивает объект App, включая настройки http сервера
// ошибка всегда nil
func New(ctx context.Context, ws *wSpace.WorkingSpace) (*App, error) {
	const (
		defaultHTTPServerWriteTimeout = time.Second * 15
		defaultHTTPServerReadTimeout  = time.Second * 15
	)

	app := new(App)
	app.ws = ws
	app.ctx = ctx

	//middleware := authorization.Authorization
	middleware := authorization.FakeAuthorization

	mux := http.NewServeMux()
	// Дать ответ клиенту о результатах вычисления выражений
	mux.HandleFunc("/getresult", middleware(ctx, ws.DB, getResult.GetResult(ctx, ws))) //

	// Получение нового выражения от клиента
	mux.HandleFunc("/newexpression", middleware(ctx, ws.DB, newExpression.NewExpression(ctx, ws)))

	// Регистрация нового пользователя
	mux.HandleFunc("/registration", registration.Registration(ctx, ws))

	// Аутентификация пользоватля
	mux.HandleFunc("/login", loginUser.Login(ctx, ws))

	// Дать задачу вычислителю // не используется
	// mux.HandleFunc("/gettask", getTask.GetTask(ctx, ws))

	// Получить ответ от вычислителя // не используется
	// mux.HandleFunc("/giveanswer", giveAnswer.GiveAnswer(ctx, ws))

	app.httpServer = &http.Server{
		Handler:      mux,
		Addr:         "localhost:" + configs.Port,
		WriteTimeout: defaultHTTPServerWriteTimeout,
		ReadTimeout:  defaultHTTPServerReadTimeout,
	}
	return app, nil
}

func (a *App) Run() error {
	log.Println("starting http server...")
	go func() {
		err := grpcServer.StartGRPCServer(a.ctx, a.ws)
		if err != nil {
			log.Fatalf("gRPC server was stop with err: %v", err)
		}
	}()
	err := a.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server was stop with err: %w", err)
	}
	log.Println("server was stop")

	return nil
}
func (a *App) stop(ctx context.Context) error {
	log.Println("shutdowning server...")
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server was shutdown with error: %w", err)
	}
	log.Println("server was shutdown")
	return nil
}

func (a *App) GracefulStop(serverCtx context.Context, sig <-chan os.Signal, serverStopCtx context.CancelFunc) {
	<-sig
	var timeOut = 30 * time.Second
	shutdownCtx, shutdownStopCtx := context.WithTimeout(serverCtx, timeOut)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			log.Fatal("graceful shutdown timed out... forcing exit")
		}
	}()

	err := a.stop(shutdownCtx)
	if err != nil {
		log.Fatal("graceful shutdown timed out... forcing exit")
	}
	serverStopCtx()
	shutdownStopCtx()
}
