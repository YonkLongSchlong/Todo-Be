package user

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/YonkLongSchlong/Todo-BE/packages/types"
	"github.com/YonkLongSchlong/Todo-BE/packages/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserRoute struct {
	store    types.UserStore
	uploader *manager.Uploader
}

func NewRoute(store types.UserStore, uploader *manager.Uploader) *UserRoute {
	return &UserRoute{
		store:    store,
		uploader: uploader,
	}
}

func (r *UserRoute) Routes(e *echo.Group) {
	e.PATCH("/user/:id", r.updateUserHandler)
	e.PATCH("/user/password/:id", r.updatePasswordHandler)
	e.PATCH("/user/avatar/:id", r.updateAvatarHandler)
}

func (r *UserRoute) updateUserHandler(c echo.Context) error {
	/** GET USER ID FROM PARAMS */
	id := c.Param("id")

	/** BIND PAYLOAD FROM CONTEXT */
	var payload types.UserPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** VALIDATE TOKEN
	 * If userId in token claims == to userId in param then continute
	 * Else return http error
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId != id {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unauthenticated user"))
	}

	/** FIND IF USER EXIST
	 * If user exists then continue
	 * Else return http error
	 */
	user, err := r.store.GetUserById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	/** COMPARE HASHED PASSWORD  AND PAYLOAD PASSWORD
	 * If password is correct then continue
	 * Else return http error
	 */
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** UPDATE USER */
	returnUser, err := r.store.UpdateUser(userId, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	returnUser.Password = ""
	return c.JSON(http.StatusOK, returnUser)
}

func (r *UserRoute) updatePasswordHandler(c echo.Context) error {
	/** GET USER ID FROM PARAMS */
	id := c.Param("id")

	/** BIND PAYLOAD FROM CONTEXT */
	var payload types.PasswordPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** VALIDATE TOKEN
	 * If userId in token claims == to userId in param then continute
	 * Else return http error
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId != id {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unauthenticated user"))
	}

	/** GET USER BY PARAMS ID TO CHECK IF USER EXIST
	 * If exits then continue
	 * Elase return http error
	 */
	user, err := r.store.GetUserById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("user don't exist"))
	}

	/** CHECK IF USER CURRENT PASSWORD MATCH
	 * If match then continue
	 * Else return http erro
	 */
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.CurrentPassword))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("password mismatch"))
	}

	/** HASH NEW PASSWORD  */
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** UPDATE PASSWORD */
	err = r.store.UpdatePassword(id, string(hashedPassword))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "update password successfully")
}

func (r *UserRoute) updateAvatarHandler(c echo.Context) error {
	/** GET USER ID FROM PARAMS */
	id := c.Param("id")

	/** GET IMAGE FILE FORM FORM DATA */
	file, err := c.FormFile("image")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** VALIDATE TOKEN
	 * If userId in token claims == to userId in param then continute
	 * Else return http error
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if userId != id {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unauthenticated user"))
	}

	/** UPDALOAD IMAGE */
	f, err := file.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer f.Close()

	result, err := r.uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("BUCKET_NAME")),
		Key:    aws.String(file.Filename),
		Body:   f,
		ACL:    "public-read",
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	/** UPDATE AVATAR IN DB */
	user, err := r.store.UpdateAvatar(id, result.Location)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)

}
