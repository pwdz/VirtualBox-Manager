package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	// "log"
	"io/ioutil"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo/v4"
	vbox "github.com/pwdz/cloudComputing/vboxWrapper"
)
const(
	CMDStatus	= "status"
	CMDOn		= "on"
	CMDOnOff	= "off"
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
	e.Any("/", endPointHandler)
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

	}else if strings.Contains(headerContentType, "multipart/form-data"){
		c.Request().ParseMultipartForm(10 << 20)
		file, handler, err := c.Request().FormFile("file")
		if err != nil{
			return err
		}
		defer file.Close()
			
		emptyFile, err := os.Create(handler.Filename)
		if err != nil {
			return err
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}
		emptyFile.Write(fileBytes)
		emptyFile.Close()

		vmName := c.Request().FormValue("vmName")
		dstPath := c.Request().FormValue("destPath") 

		cmd = command{
			Type: CMDUpload,
			DestVM: vmName,
			DestPath: dstPath,
			OriginPath: handler.Filename,
		}
	}

	handleCommand(cmd)

	return c.String(http.StatusOK, "")
}
func handleCommand(cmd command){
	switch cmd.Type{
	case CMDStatus:
		handleStatus(cmd)
	case CMDOn, CMDOnOff:
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
	if cmd.Type == CMDOn{
		vbox.PowerOn(cmd.VmName)
	}else{
		vbox.PowerOff(cmd.VmName)
	}
}
func handleDelete(cmd command){
	vbox.Delete(cmd.VmName)
}
func handleSetting(cmd command){
	vbox.ChangeSetting(cmd.VmName, cmd.Cpu, cmd.Ram)
}
func handleTransfer(cmd command){
	vbox.Transfer(cmd.OriginVM, cmd.DestVM, cmd.OriginPath, cmd.DestPath)
}
func handleClone(cmd command){
	vbox.Clone(cmd.SourceVmName, cmd.DestVmName)
}
func handleExecute(cmd command){
	vbox.Execute(cmd.VmName, cmd.Input)
}
func handleUpload(cmd command){
	fmt.Println(cmd)
	vbox.Upload(cmd.DestVM, cmd.DestPath, cmd.OriginPath)
}