# VirtualBox-Manager
A simple http server and virtualbox wrapper for controlling virtualbox remotly.  
requests must be sent like the below json format:  
```json
{
  "command": "delete",
  "vmName": "VM1",
  "status": "Ok"
}
```
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
