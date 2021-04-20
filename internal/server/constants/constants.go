package constants
import  "github.com/pwdz/cloudComputing/internal/server/user"
const(
	CMDStatus	= "status"
	CMDOn		= "on"
	CMDOff	= "off"
	CMDSetting	= "setting"
	CMDClone	= "clone"
	CMDDelete	= "delete"
	CMDExecute	= "execute"
	CMDTransfer	= "transfer"
	CMDUpload	= "upload"
)
// Create the JWT key used to create the signature
var JwtKey = []byte("vboxmanagerKey")

var Users = []user.User{
	user.NewUser("admin", "adminpass", "Admin"),
	user.NewUser("user","userpass","User"),
}