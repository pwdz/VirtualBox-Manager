package handler

import(
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors"
	"github.com/pwdz/cloudComputing/internal/server/constants"
	"github.com/pwdz/cloudComputing/internal/server/user"
	"os"
	"io/ioutil"
	"encoding/json"
	"fmt"
)

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create the Signin handler
func Login(c echo.Context) error{
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(c.Request().Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		return c.String(http.StatusBadRequest, "")
	}

	// Get the expected password from our in memory map
	ok := false
	var expectedPassword, role string
	for _, user := range constants.Users{ 
		if user.Username == creds.Username{
			expectedPassword = user.Password
			ok = true
			role = user.Role
		}
	}

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		return c.String(http.StatusUnauthorized, "")
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	// expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &user.User{
		Username: creds.Username,
		Role: role,
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(constants.JwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return c.String(http.StatusInternalServerError, "")
	}

	return c.String(http.StatusOK, tokenString)
}

func EndPointHandler(c echo.Context) error{
	
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
			Type: constants.CMDUpload,
			DestVM: vmName,
			DestPath: dstPath,
			OriginPath: handler.Filename,
		}
	}
	role := fmt.Sprintf("%v", c.Get("role"))
	jsonResponse := cmd.handleCommand(role)

	return c.JSONBlob(http.StatusOK, jsonResponse)
}
