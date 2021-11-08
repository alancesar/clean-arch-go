package api

import (
	"clean-arch/date"
	"clean-arch/todo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userKey = "x-user"

type Reader interface {
	ListAllPending(user string) ([]todo.Todo, error)
}

type Writer interface {
	MarkAsDoneByUserAndID(user string, id int) (todo.Todo, error)
}

type markAsUpdateRequest struct {
	ID int `uri:"id"`
}

type response struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	DueDate string `json:"due_date,omitempty"`
}

func ListPendingTodosHandler(reader Reader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		u, err := retrieveUserHeader(ctx)
		if err != nil {
			badRequestError(ctx, err)
			return
		}

		if todos, err := reader.ListAllPending(u); err != nil {
			internalServerError(ctx)
		} else {
			r := parseResponseFromEntities(todos)
			ctx.JSON(http.StatusOK, r)
		}
	}
}

func MarkTodoAsDoneHandler(writer Writer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request markAsUpdateRequest
		if err := ctx.ShouldBindUri(&request); err != nil {
			badRequestError(ctx, err)
			return
		}

		u, err := retrieveUserHeader(ctx)
		if err != nil {
			badRequestError(ctx, err)
			return
		}

		if t, err := writer.MarkAsDoneByUserAndID(u, request.ID); err != nil {
			handleError(ctx, err)
		} else {
			ctx.JSON(http.StatusOK, t)
		}
	}
}

func parseResponseFromEntities(entities []todo.Todo) []response {
	output := make([]response, len(entities))
	for index, entity := range entities {
		output[index] = response{
			ID:    entity.ID,
			Title: entity.Title,
		}

		if entity.DueDate != nil {
			output[index].DueDate = date.ParseToString(*entity.DueDate)
		}
	}

	return output
}

func handleError(ctx *gin.Context, err error) {
	if err == todo.AlreadyDoneErr {
		badRequestError(ctx, err)
		return
	}

	internalServerError(ctx)
}

func badRequestError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func internalServerError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "unexpected error, please try again later",
	})
}

func retrieveUserHeader(ctx *gin.Context) (string, error) {
	if u := ctx.GetHeader(userKey); u == "" {
		return "", errors.New("header 'x-user' must be provided")
	} else {
		return u, nil
	}
}
