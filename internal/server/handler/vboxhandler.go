package handler
import(
	tools "github.com/pwdz/cloudComputing/pkg"
	vbox "github.com/pwdz/cloudComputing/vboxWrapper"
	"github.com/pwdz/cloudComputing/internal/server/constants"
	"encoding/json"
)


type command struct{
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

type response struct{
	Err			string				`json:"error,omitempty"`
	Status		string				`json:"status,omitempty"`
	Response	string				`json:"response,omitempty"`
	Details		[]map[string]string	`json:"details,omitempty"`
}

func (cmd command) handleCommand(role string) []byte{
	if role == "User"{
		if 	cmd.VmName != "" && cmd.VmName != "vm1" ||
			cmd.OriginVM != "" && cmd.OriginVM != "vm1" ||
			cmd.DestVM != "" && cmd.DestVM != "vm1" {
				return []byte("{\"error\": \"You don't have access to this vm\"}")
			}

	}
	switch cmd.Type{
	case constants.CMDStatus:
		return cmd.handleStatus()
	case constants.CMDOn, constants.CMDOff:
		return cmd.handleOnOff()
	case constants.CMDDelete:
		return cmd.handleDelete()
	case constants.CMDSetting:
		return cmd.handleSetting()
	case constants.CMDTransfer:
		return cmd.handleTransfer()
	case constants.CMDClone:
		return cmd.handleClone()
	case constants.CMDExecute:
		return cmd.handleExecute()
	case constants.CMDUpload:
		return cmd.handleUpload()
	}
	return []byte("{\"error\": \"Invalid command\"}")
}

func (cmd command) handleStatus() []byte{
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
func (cmd command) handleOnOff() []byte{
	var status string
	var err error
	if cmd.Type == constants.CMDOn{
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
func (cmd command) handleDelete() []byte{
	status, err := vbox.Delete(cmd.VmName)
	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)

	return tools.ConcatJsons(cmdJson, respJson)
}
func (cmd command) handleSetting() []byte{
	status, err := vbox.ChangeSetting(cmd.VmName, cmd.Cpu, cmd.Ram)

	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)

	return tools.ConcatJsons(cmdJson, respJson)
}
func (cmd command) handleTransfer()[]byte{
	status, err := vbox.Transfer(cmd.OriginVM, cmd.DestVM, cmd.OriginPath, cmd.DestPath)
	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)

	return tools.ConcatJsons(cmdJson, respJson)
}
func (cmd command) handleClone() []byte{
	status, err := vbox.Clone(cmd.SourceVmName, cmd.DestVmName)
	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)

	return tools.ConcatJsons(cmdJson, respJson)
}
func (cmd command) handleExecute()[]byte{
	status, rsp, err := vbox.Execute(cmd.VmName, cmd.Input)
	resp := response{Status: status, Response: rsp}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)

	return tools.ConcatJsons(cmdJson, respJson)
}
func (cmd command) handleUpload() []byte{
	status, err := vbox.Upload(cmd.DestVM, cmd.DestPath, cmd.OriginPath)
	
	resp := response{Status: status}
	if err != nil{
		resp.Err = err.Error()
	}

	cmdJson, _ := json.Marshal(cmd)
	respJson, _ := json.Marshal(resp)

	return tools.ConcatJsons(cmdJson, respJson)
}