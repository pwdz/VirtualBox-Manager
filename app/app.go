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
	tools "github.com/pwdz/cloudComputing/pkg"
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

type( 
	command struct{
		Type 			string 	`json:"command,omitempty" required:"true"`
		VmName			string	`json:"vmName,omitempty"`
		Cpu 			int		`json:"cpu,omitempty"`
		Ram 			int		`json:"ram,omitempty"`
		SourceVmName	string	`json:"sourceVmName,omitempty"`
		DestVmName		string  `json:"destVmName,omitempty"`
		Input			string	`json:"input,omitempty"`
		OriginVM		string	`json:"originVM,omitempty"`
		OriginPath		string	`json:"originPath,omitempty"`
		DestVM			string	`json:"destVM,omitempty"`
		DestPath		string	`json:"destPath,omitempty"`
	}

	response struct{
		Err			string				`json:"error,omitempty"`
		Status		string				`json:"status,omitempty"`
		Response	string				`json:"response,omitempty"`
		Details		[]map[string]string	`json:"details,omitempty"`
	}
)

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

// func middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
// 		if len(authHeader) != 2 {
// 			fmt.Println("Malformed token")
// 			w.WriteHeader(http.StatusUnauthorized)
// 			w.Write([]byte("Malformed Token"))
// 		} else {
// 			jwtToken := authHeader[1]
// 			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (string{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 				}
// 				return []byte(SECRETKEY), nil
// 			})

// 			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 				ctx := context.WithValue(r.Context(), "props", claims)
// 				// Access context values in handlers like this
// 				// props, _ := r.Context().Value("props").(jwt.MapClaims)
// 				next.ServeHTTP(w, r.WithContext(ctx))
// 			} else {
// 				fmt.Println(err)
// 				w.WriteHeader(http.StatusUnauthorized)
// 				w.Write([]byte("Unauthorized"))
// 			}
// 		}
// 	})
// }
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
		bytess := handleStatus(cmd)
		log.Println(string(bytess))
	case CMDOn, CMDOnOff:
		bytess:=handleOnOff(cmd)
		log.Println(string(bytess))
	case CMDDelete:
		bytess:=handleDelete(cmd)
		log.Println(string(bytess))
	case CMDSetting:
		handleSetting(cmd)
	case CMDTransfer:
		handleTransfer(cmd)
	case CMDClone:
		bytess:=handleClone(cmd)
		log.Println(string(bytess))
	case CMDExecute:
		handleExecute(cmd)
	case CMDUpload:
		handleUpload(cmd)
	}
}
func handleStatus(cmd command) []byte{
	resp := response{}
	var respJson []byte

	if cmd.VmName != ""{
		status, err := vbox.GetStatus(cmd.VmName)
		resp.Status = status
		if err != nil{
			resp.Err = err.Error()
		}
	}else{
		vmNames, err := vbox.GetVmNames()

		if err != nil{
			resp.Err = err.Error()							
		}else{
			statuses := make([]map[string]string, len(vmNames))
			for index, vmName := range vmNames{
				statuses[index] = map[string]string{"vmName": vmName}

				status, err := vbox.GetStatus(vmName)
				statuses[index]["vmName"] = vmName
				if err != nil{
					statuses[index]["error"] = err.Error()
				}else{
					statuses[index]["status"] = status
				}
			} 
			resp.Details = statuses
		}
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ = json.Marshal(resp)
	return tools.ConcatJsons(cmdJson, respJson)
}
func handleOnOff(cmd command) []byte{
	var status string
	var err error
	if cmd.Type == CMDOn{
		status, err = vbox.PowerOn(cmd.VmName)
	}else{
		status, err = vbox.PowerOff(cmd.VmName)
	}
	cmdJson, _ := json.Marshal(cmd)

	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}
	respJson, _ := json.Marshal(resp)
	return tools.ConcatJsons(cmdJson, respJson)
}
func handleDelete(cmd command)[]byte{
	status, err := vbox.Delete(cmd.VmName)
	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)
	return tools.ConcatJsons(cmdJson, respJson)
}
func handleSetting(cmd command){
	vbox.ChangeSetting(cmd.VmName, cmd.Cpu, cmd.Ram)
}
func handleTransfer(cmd command){
	vbox.Transfer(cmd.OriginVM, cmd.DestVM, cmd.OriginPath, cmd.DestPath)
}
func handleClone(cmd command)[]byte{
	status, err := vbox.Clone(cmd.SourceVmName, cmd.DestVmName)
	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)
	return tools.ConcatJsons(cmdJson, respJson)

}
func handleExecute(cmd command){
	vbox.Execute(cmd.VmName, cmd.Input)
}
func handleUpload(cmd command){
	fmt.Println(cmd)
	vbox.Upload(cmd.DestVM, cmd.DestPath, cmd.OriginPath)
}