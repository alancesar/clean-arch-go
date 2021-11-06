package api

import (
	"clean-arch/date"
	"clean-arch/todo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Reader interface {
	ListAllPending() ([]todo.Todo, error)
}

type Writer interface {
	MarkAsDoneByID(id int) (todo.Todo, error)
}

type markAsUpdateRequest struct {
	ID int `uri:"id"`
}

type response struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	DueDate string `json:"due_date"`
}

func ListPendingTodosHandler(reader Reader) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if todos, err := reader.ListAllPending(); err != nil {
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

		if t, err := writer.MarkAsDoneByID(request.ID); err != nil {
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
			ID:      entity.ID,
			Title:   entity.Title,
			DueDate: date.ParseToString(entity.DueDate),
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
