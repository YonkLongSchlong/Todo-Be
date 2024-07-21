package auth

import (
	"fmt"
	"net/http"

	"github.com/YonkLongSchlong/Todo-BE/packages/types"
	"github.com/YonkLongSchlong/Todo-BE/packages/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthRoute struct {
	store types.UserStore
}

func NewRoute(store types.UserStore) *AuthRoute {
	return &AuthRoute{store: store}
}

func (r *AuthRoute) Routes(e *echo.Group) {
	e.POST("/register", r.RegisterHanlder)
	e.POST("/login", r.LoginHandler)
}

func (r *AuthRoute) RegisterHanlder(c echo.Context) error {
	/** BIND PAYLOAD FROM CONTEXT */
	var payload types.RegisterPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** CREATE USER */
	err = r.store.CreateUser(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "Create user successfully")
}

func (r *AuthRoute) LoginHandler(c echo.Context) error {
	/** BIND PAYLOAD FROM CONTEXT */
	var payload types.LoginPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** GET USER BY EMAIL
	 * If user exist continue
	 * Else return http error
	 */
	user, err := r.store.GetUserByEmail(payload.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** COMPARE HASH PASSWORD TO PAYLOAD INPUT PASSWORD
	 * If check == true continue
	 * Else return http error
	 */
	check := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if check != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("invalid credentials, please try again"))
	}

	/** GENERATE JWT TOKEN */
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** RETURN USER JSON
	 * Set user password to an empty string first before returning it
	 */
	user.Password = ""
	return c.JSON(http.StatusOK, map[string]any{
		"user":  user,
		"token": token,
	})
}
