package todo

import (
	"fmt"
	"net/http"

	"github.com/YonkLongSchlong/Todo-BE/packages/types"
	"github.com/YonkLongSchlong/Todo-BE/packages/utils"
	"github.com/labstack/echo/v4"
)

type TodoRoute struct {
	store types.TodoStore
}

func NewRoute(store types.TodoStore) *TodoRoute {
	return &TodoRoute{store: store}
}

func (r *TodoRoute) Routes(e *echo.Group) {
	e.POST("/todo/create", r.createTodoHandler)
	e.DELETE("/todo/:id", r.deleteTodoHandler)
	e.PATCH("/todo/:id", r.updpateTodoHandler)
	e.PATCH("/todo/status/:id", r.updateTodoStatusHandler)
	e.GET("/todo/:id", r.GetTodoByIdHandler)
	e.POST("/todo", r.GetTodoByDate)
}

func (r *TodoRoute) createTodoHandler(c echo.Context) error {
	/** BIND TODO PAYLOAD FROM JSON */
	var payload types.TodoPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** VALIDATE TOKEN
	 * If userId in token claims == to userId in Payload then continute
	 * Else return http error
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if userId != payload.UserId {
		return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unauthenticated user"))
	}

	/** CREATE TODO */
	err = r.store.CreateTodo(payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, "create todo successfully")
}

func (r *TodoRoute) deleteTodoHandler(c echo.Context) error {
	/** GET TODO ID FROM PARAMS */
	todoId := c.Param("id")
	if todoId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing id parameter"))
	}

	/** VALIDATE TOKEN
	 * Get the userId from token claims to user it for querying purposes
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** DELETE TODO */
	err = r.store.DeleteTodo(todoId, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, "delete todo successfully")
}

func (r *TodoRoute) updpateTodoHandler(c echo.Context) error {
	/** GET TODO ID FROM PARAMS */
	todoId := c.Param("id")
	if todoId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing id parameter"))
	}

	/** BIND TODO PAYLOAD FROM JSON */
	var payload types.TodoPayload
	err := c.Bind(&payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** VALIDATE TOKEN
	 * Get the userId from token claims to user it for querying purposes
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** UPDATE TODO */
	err = r.store.UpdateTodo(todoId, userId, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "updated successfully")
}

func (r *TodoRoute) updateTodoStatusHandler(c echo.Context) error {
	/** GET TODO ID FROM PARAMS */
	todoId := c.Param("id")
	if todoId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing id parameter"))
	}

	/** VALIDATE TOKEN
	 * Get the userId from token claims to user it for querying purposes
	 */
	userId, err := utils.GetClaims(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	/** UPDATE IS_COMPLETED */
	err = r.store.SetIsCompledtedTodo(todoId, userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "updated status successfully")
}

func (r *TodoRoute) GetTodoByIdHandler(c echo.Context) error {
	/** GET TODO ID FROM PARAMS */
	todoId := c.Param("id")
	if todoId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("missing id parameter"))
	}

	/** GET TODO BY ID */
	todo, err := r.store.GetTodoById(todoId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())

	}

	return c.JSON(http.StatusOK, todo)
}

func (r *TodoRoute) GetTodoByDate(c echo.Context) error {
	var jsonMap map[string]string
	err := c.Bind(&jsonMap)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	todo, err := r.store.GetTodoByDate(jsonMap["date"])
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, todo)
}
