# VirtualBox-Manager
A simple http server(using [echo](https://github.com/labstack/echo)) and virtualbox wrapper for controlling virtualbox vms remotly.  

available commands:  
```
status
on
off
setting
clone
delete
execute
transfer
upload
```
More details:  
```json
{
  "command": "status",
  "vmName": "VM1"
}
```
```json
{
"command": "status"
}
```
```json
{
  "command": "delete",
  "vmName": "VM1",
  "status": "Ok"
}
```
```json
{
  "command": "on/off",
  "vmName": "VM1"
}
```
```json
{
  "command": "setting",
  "vmName": "VM1",
  "cpu": 2,
  "ram": 1024
}
```
```json
{
  "command": "clone",
  "sourceVmName": "VM1",
  "destVmName": "VM2"
}
```
```json
{
  "command": "execute",
  "vmName": "VM1",
  "input": "mkdir sina && touch sina.txt && ls"
}
```
```json
{
  "command": "transfer",
  "originVM": "VM1",
  "originPath": "/home/sina.txt",
  "destVM": "VM2",
  "destPath": "/home/temp/"
}
```
Upload file in `multipart/form-data` format.
```
destPath
vmName
file
```
