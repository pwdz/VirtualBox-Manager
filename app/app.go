package app

import (
	"fmt"
	// "log"
	// "log"
	// "io/ioutil"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	vbox "github.com/pwdz/cloudComputing/vboxWrapper"
)
const(
	CMDStatus	= "status"
	CMDOnOff	= "on/off"
	CMDSetting	= "setting"
	CMDClone	= "clone"
	CMDDelete	= "delete"
	CMDExecute	= "execute"
	CMDTransfer	= "transfer"
	CMDUpload	= "upload"
)
var e* echo.Echo

type command struct{
	Type 			string 	`json:"command"`
	VmName			string	`json:"vmName"`
	Cpu 			int		`json:"cpu"`
	Ram 			int		`json:"ram"`
	SourceVmName	string	`json:"sourceVmName"`
	DestVmName		string  `json:"destVmName"`
	Input			string	`json:"input"`
	OriginVM		string	`json:"originVM"`
	OriginPath		string	`json:"originPath"`
	DestVM			string	`json:"destVM"`
	DestPath		string	`json:"destPath"`
}

func InitCfg(){
	err := cleanenv.ReadEnv(&Cfg)
	fmt.Printf("%+v", Cfg)
	if err != nil{
		e.Logger.Fatal("Unable to load configs")
	}
}
func InitServer(){
	e = echo.New()
	e.GET("/", endPointHandler)
	e.Logger.Fatal(e.Start(Cfg.Host + ":" + Cfg.Port))
}
func endPointHandler(c echo.Context) error{
	headerContentType := c.Request().Header.Get("Content-Type")
	var cmd command

	if headerContentType == "application/json" {		
		var unmarshalErr *json.UnmarshalTypeError

		decoder := json.NewDecoder(c.Request().Body)
		decoder.DisallowUnknownFields()

		err := decoder.Decode(&cmd)
		if err != nil {
			if errors.As(err, &unmarshalErr) {
				return c.String(http.StatusBadRequest, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field)
			} else {
				return c.String(http.StatusBadRequest, "Bad Request "+err.Error())
			}
		}

		handleCommand(cmd)
	}else{
		//TODO: upload file
	}

	return c.String(http.StatusOK, "")
}
func handleCommand(cmd command){
	switch cmd.Type{
	case CMDStatus:
		handleStatus(cmd)
	case CMDOnOff:
		handleOnOff(cmd)
	case CMDDelete:
		handleDelete(cmd)
	case CMDSetting:
		handleSetting(cmd)
	case CMDTransfer:
		handleTransfer(cmd)
	case CMDClone:
		handleClone(cmd)
	case CMDExecute:
		handleExecute(cmd)
	case CMDUpload:
		handleUpload(cmd)
	}
}
func handleStatus(cmd command) string{
	status, err := vbox.GetStatus(cmd.VmName)
	if err != nil{
		
	}
	fmt.Println(status)
	return status
}
func handleOnOff(cmd command){

}
func handleDelete(cmd command){

}
func handleSetting(cmd command){

}
func handleTransfer(cmd command){

}
func handleClone(cmd command){

}
func handleExecute(cmd command){

}
func handleUpload(cmd command){
	
}