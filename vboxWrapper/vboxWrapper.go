package vboxWrapper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

/*
	vboxmanage list vms
	vboxmanage list runningvms
	vboxmanage startvm <name or UUID>
	vboxmanage controlvm <subcommand>
		pause, resume, reset, poweroff, and savestate
	vboxmanage unregister <name or UUID> --delete
	vboxmanage showvminfo <name or UUID>

	vboxmanage modifyvm <name or UUID> --memory <RAM in MB>
	vboxmanage modifyvm <name or UUID> --cpus <number>
*/
const(
	VBoxCommand = "vboxmanage"
)

func printCommand(cmd *exec.Cmd) {
  log.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
  if err != nil {
    os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
  }
}

func printOutput(outs []byte) {
  if len(outs) > 0 {
    log.Printf("==> Output: %s\n", string(outs))
  }
}

func GetStatus(vmName string) (string, error){
	cmd := exec.Command(VBoxCommand, "showvminfo",vmName,"--machinereadable")

	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil{
		printError(err)
		return "", err
	}
	regex, _ := regexp.Compile("VMState=\"[a-zA-Z]+\"")
	status := regex.FindString(string(output))
	log.Println(status)
	status = strings.Split(status, "=")[1]
	status = strings.Trim(status, "\"")
	return status, nil
}

func PowerOn(vmName string){
	status, _ := GetStatus(vmName)
	if status == "poweroff" {
		cmd := exec.Command(VBoxCommand, "startvm",vmName,"--type","headless")

		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printOutput(output)
		printError(err)
	}
}
func PowerOff(vmName string){
	status, _ := GetStatus(vmName)
	if status == "running" {
		cmd := exec.Command(VBoxCommand, "controlvm",vmName,"poweroff")

		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printOutput(output)
		printError(err)
	}
}

func ChangeSetting(vmName string, cpu, ram int){
	if cpu > 0{
		cmd := exec.Command(VBoxCommand, "modifyvm",vmName,"--cpus",strconv.Itoa(cpu))
		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printOutput(output)
		printError(err)
	}
	if ram > 0{
		cmd := exec.Command(VBoxCommand, "modifyvm",vmName,"--memory",strconv.Itoa(ram))
		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printOutput(output)
		printError(err)
	}
}