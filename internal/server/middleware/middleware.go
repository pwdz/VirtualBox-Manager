package middleware
import(
	"github.com/labstack/echo/v4"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"log"
	"github.com/pwdz/cloudComputing/internal/server/user"
	"github.com/pwdz/cloudComputing/internal/server/constants"
)

func Authorize(next echo.HandlerFunc)echo.HandlerFunc {
	return func(c echo.Context) error {
		log.Println("Middleware")
		// Get the JWT string from the header
		tknStr := c.Request().Header.Get("token")

		// Initialize a new instance of `Claims`
		claims := &user.User{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return constants.JwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				return c.String(http.StatusUnauthorized, "Unauthorized")
			}
			return c.String(http.StatusBadRequest, "Bad request")
		}
		if !tkn.Valid {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		// Finally, return the welcome message to the user, along with their
		// username given in the token
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		return next(c)
	}
}

