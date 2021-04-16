package vboxWrapper

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
	PowerOn = iota
	PowerOff
)

func printCommand(cmd *exec.Cmd) {
  fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
  if err != nil {
    os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
  }
}

func printOutput(outs []byte) {
  if len(outs) > 0 {
    fmt.Printf("==> Output: %s\n", string(outs))
  }
}

func GetStatus(vmName string) (string, error){
	cmd := exec.Command("vboxmanage", "showvminfo",vmName,"--machinereadable")

	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	if err != nil{
		printError(err)
		return "", err
	}
	regex, _ := regexp.Compile("VMState=\"[a-zA-Z]+\"")
	status := regex.FindString(string(output))
	return strings.Split(status, "=")[1], nil
}

func PowerOffNo(status int){
	if status == PowerOff{

	}else if status == PowerOn{

	}
}