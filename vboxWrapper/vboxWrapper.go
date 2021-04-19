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

func GetStatus(vmName string) string{
	cmd := exec.Command(VBoxCommand, "showvminfo",vmName,"--machinereadable")

	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	if err != nil{
		printError(err)
		return err.Error()
	}
	regex, _ := regexp.Compile("VMState=\"[a-zA-Z]+\"")
	status := regex.FindString(string(output))
	log.Println(status)
	status = strings.Split(status, "=")[1]
	status = strings.Trim(status, "\"")
	return status
}
func GetVmNames() ([]string, error){
	cmd := exec.Command(VBoxCommand, "list","vms")

	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	
	if err != nil{
		printError(err)
		return nil, err
	}
	regex, _ := regexp.Compile("\"[A-Za-z0-9]+\"")
	vmNames := regex.FindAllString(string(output), -1)

	for index, vmName := range vmNames{
		vmNames[index] = strings.Trim(vmName, "\"")
	}
	return vmNames, nil
}
func PowerOn(vmName string) string{
	status := GetStatus(vmName)
	if status == "poweroff" {
		cmd := exec.Command(VBoxCommand, "startvm",vmName,"--type","headless")

		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		if err != nil{
			printError(err)
			return err.Error()
		}
		printOutput(output)
		return ""
	}
	return ""
}
func PowerOff(vmName string){
	status := GetStatus(vmName)
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
func Clone(vmSrc, vmDst string){
	cmd := exec.Command(VBoxCommand, "clonevm",vmSrc,"--name",vmDst, "--register")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	printError(err)	
}
func Delete(vmName string){
	cmd := exec.Command(VBoxCommand, "unregistervm",vmName,"--delete")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	printError(err)	
}
func Execute(vmName, input string){
	cmd := exec.Command(VBoxCommand, "guestcontrol",vmName,"run","bin/sh","--username","pwdz","--password", "pwdz", "--wait-stdout", "--wait-stderr", "--","-c",input)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	printError(err)	
}
func Transfer(vmSrc, vmDst, originPath, dstPath string){
	 _, err := os.Stat("./temp")
    if err != nil && os.IsNotExist(err){
		os.Mkdir("temp", 666)
	}

	paths := strings.Split(originPath, "/")
	fileName := paths[len(paths) - 1]

	internalPath :=   "./temp/" 
	log.Println(vmSrc, vmDst)

	copyFromCommand := exec.Command(VBoxCommand, "guestcontrol",vmSrc,"copyfrom","--target-directory",internalPath , originPath,"--username","pwdz","--password", "pwdz")
	printCommand(copyFromCommand)
	copyFromOutput, err := copyFromCommand.CombinedOutput()
	printOutput(copyFromOutput)
	printError(err)	

	internalPath += fileName
	copyToCommand := exec.Command(VBoxCommand, "guestcontrol",vmDst,"copyto","--target-directory",dstPath, internalPath,"--username","pwdz","--password", "pwdz")
	printCommand(copyToCommand)
	copyToOutput, err := copyToCommand.CombinedOutput()
	printOutput(copyToOutput)
	printError(err)	
}
func Upload(vmDst, dstPath, originPath string ){
	copyToCommand := exec.Command(VBoxCommand, "guestcontrol",vmDst,"copyto","--target-directory",dstPath, "D:/AUT/Courses/Term6/Cloud Computing/CloudComputing/TestFile.txt","--username","pwdz","--password", "pwdz", "--verbose")
	printCommand(copyToCommand)
	copyToOutput, err := copyToCommand.CombinedOutput()
	printOutput(copyToOutput)
	printError(err)	
}