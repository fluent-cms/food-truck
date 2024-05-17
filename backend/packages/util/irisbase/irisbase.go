package irisbase

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/cors"
	irisLogger "github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"github.com/kataras/iris/v12/mvc"
	"github.com/samber/lo"
)

type AppConfig struct {
	ApiPrefix string `yaml:"apiPrefix"`
	Debug     bool   `yaml:"debug"`
	Port      int    `yaml:"port"`
}

type errHandler func(ctx iris.Context, err error)
type App struct {
	Config     AppConfig
	IrisApp    *iris.Application
	ErrHandler errHandler
}

type AppBuilder interface {
	BadRequest() []error
	Services() []any
	Controller() map[string]any
}

func createErrorHandler(isDebug bool, badRequests []error) errHandler {
	return func(ctx iris.Context, err error) {
		id := uuid.New().ID()
		code := 500
		for _, err2 := range badRequests {
			if errors.As(err, &err2) {
				code = 400
				break
			}
		}

		if isDebug {
			ctx.StopWithError(code, fmt.Errorf("%w, id=%v", err, id))
		} else {
			if code == 400 {
				ctx.StopWithError(code, fmt.Errorf("bad request, id=%v", id))
			} else {
				ctx.StopWithError(code, fmt.Errorf("internal error, id=%v", id))
			}
		}
		fmt.Printf("API-ERR: %v, ErrID=%v\n", err, id)
	}
}

func NewIrisApp(config AppConfig, builder AppBuilder) *App {
	i := &App{
		Config: config,
	}
	i.ErrHandler = createErrorHandler(config.Debug, builder.BadRequest())
	app := iris.New()
	app.Use(recover.New())
	app.Use(irisLogger.New())
	if i.Config.Debug {
		fmt.Println("************************")
		fmt.Println("running on debug mode...")
		fmt.Println("************************")
		app.Logger().SetLevel("debug")
		app.UseRouter(cors.New().
			HandleErrorFunc(i.ErrHandler).
			ExtractOriginFunc(cors.DefaultOriginExtractor).
			ReferrerPolicy(cors.NoReferrerWhenDowngrade).
			AllowOrigin("*").
			AllowHeaders("content-type, authorization").
			Handler())
	}
	app.RegisterDependency(lo.ToAnySlice(builder.Services())...)
	for s, a := range builder.Controller() {
		mvc.New(app.Party(i.Config.ApiPrefix + s)).Handle(a).HandleError(i.ErrHandler)
	}
	app.HandleDir("/", iris.Dir("./web"), iris.DirOptions{
		IndexName: "index.html",
		SPA:       true,
	})
	i.IrisApp = app
	return i
}

func (i *App) AddRouter(relPath string, handler any) {
	mvc.New(i.IrisApp.Party(i.Config.ApiPrefix + relPath)).Handle(handler).HandleError(i.ErrHandler)
}

func (i *App) Start() {
	if err := i.IrisApp.Listen(fmt.Sprintf(":%d", i.Config.Port), func(app *iris.Application) {
		app.Configure(
			iris.WithOptimizations,
			iris.WithFireMethodNotAllowed,
			iris.WithPathIntelligence,
		)
	}); err != nil {
		panic(err)
	}
}
