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
func PowerOn(vmName string) (string, error){
	status, err := GetStatus(vmName)
	if err != nil{
		return "", err
	}
	
	if status == "poweroff" {
		cmd := exec.Command(VBoxCommand, "startvm",vmName,"--type","headless")

		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printOutput(output)
		if err != nil{
			printError(err)
			return "", err
		}
		return "Powering on", nil
	}
	return "", fmt.Errorf(vmName + ">> current status: " + status)
}
func PowerOff(vmName string)(string, error){
	status, err := GetStatus(vmName)
	if err != nil{
		return "", err
	}
	
	if status == "running" {
		cmd := exec.Command(VBoxCommand, "controlvm",vmName,"poweroff")

		printCommand(cmd)
		output, err := cmd.CombinedOutput()
		printOutput(output)
		printError(err)
		if err != nil{
			printError(err)
			return "", err
		}
		return "Powering off", nil
	}

	return "", fmt.Errorf(vmName + ">> current status: " + status)
}

func ChangeSetting(vmName string, cpu, ram int)(string, error){
	args := []string{"modifyvm", vmName}

	if cpu > 0{
		args = append(args, "--cpus", strconv.Itoa(cpu))
	}
	if ram > 0{
		args = append(args, "--memory",strconv.Itoa(ram))
	}
	
	cmd := exec.Command(VBoxCommand, args...)	
	printCommand(cmd)

	output, err := cmd.CombinedOutput()
	printOutput(output)
	if err != nil{
		printError(err)
		return "", fmt.Errorf(string(output))
	}

	return "Ok", nil
}
func Clone(vmSrc, vmDst string)(string, error){
	cmd := exec.Command(VBoxCommand, "clonevm",vmSrc,"--name",vmDst, "--register")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	if err != nil{
		printError(err)
		return "", fmt.Errorf(string(output))
	}

	return "Ok", nil
}
func Delete(vmName string)(string, error){
	cmd := exec.Command(VBoxCommand, "unregistervm",vmName,"--delete")
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	if err != nil{
		printError(err)
		return "", fmt.Errorf(string(output))
	}

	return "Ok", nil
}
func Execute(vmName, input string)(string, string, error){
	cmd := exec.Command(VBoxCommand, "guestcontrol",vmName,"run","bin/sh","--username","pwdz","--password", "pwdz", "--wait-stdout", "--wait-stderr", "--","-c",input)
	printCommand(cmd)
	output, err := cmd.CombinedOutput()
	printOutput(output)
	if err != nil{
		printError(err)
		return "", "", fmt.Errorf(string(output))
	}

	return "Ok", string(output), nil
}
func Transfer(vmSrc, vmDst, originPath, dstPath string)(string, error){
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
	if err != nil{
		printError(err)
		return "", fmt.Errorf(string(copyFromOutput))
	}


	internalPath += fileName
	copyToCommand := exec.Command(VBoxCommand, "guestcontrol",vmDst,"copyto","--target-directory",dstPath, internalPath,"--username","pwdz","--password", "pwdz")
	printCommand(copyToCommand)
	copyToOutput, err := copyToCommand.CombinedOutput()
	printOutput(copyToOutput)
	if err != nil{
		printError(err)
		return "", fmt.Errorf(string(copyToOutput))
	}

	return "Ok", nil
}
func Upload(vmDst, dstPath, originPath string)(string, error){
	copyToCommand := exec.Command(VBoxCommand, "guestcontrol",vmDst,"copyto","--target-directory",dstPath, "D:/AUT/Courses/Term6/Cloud Computing/CloudComputing/TestFile.txt","--username","pwdz","--password", "pwdz", "--verbose")
	printCommand(copyToCommand)
	copyToOutput, err := copyToCommand.CombinedOutput()
	printOutput(copyToOutput)
	if err != nil{
		printError(err)
		return "", fmt.Errorf(string(copyToOutput))
	}

	return "Ok", nil
}