package server

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	"github.com/pwdz/cloudComputing/internal/server/handler"
	"github.com/pwdz/cloudComputing/internal/server/middleware"
	"log"
)

var e* echo.Echo

func InitCfg(){
	err := cleanenv.ReadEnv(&Cfg)
	log.Printf("%+v", Cfg)
	if err != nil{
		e.Logger.Fatal("Unable to load configs")
	}
}
func InitServer(){
	e = echo.New()
	e.Any("/", handler.EndPointHandler, middleware.Authorize)
	e.GET("/login", handler.Login)
	e.Logger.Fatal(e.Start(Cfg.Host + ":" + Cfg.Port))
}
