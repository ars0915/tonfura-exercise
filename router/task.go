package router

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/ars0915/gogolook-exercise/entity"
	"github.com/ars0915/gogolook-exercise/util/cGin"
)

func (rH *HttpHandler) ListTasksHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	page := ctx.GetPaginator()
	param := entity.ListTaskParam{
		Offset: &page.Offset,
		Limit:  &page.Limit,
	}

	data, count, err := rH.h.ListTasks(ctx, param)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "List Tasks Failed")
		return
	}
	page.SetTotalCount(int(count))

	ctx.WithPaginator(page).WithData(data).Response(http.StatusOK, "")
}

type createTaskBody struct {
	Name   string `json:"name" binding:"required,gt=0,lte=255"`
	Status uint8  `json:"status" binding:"required,oneof=0 1"`
}

func (rH *HttpHandler) CreateTaskHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	var body createTaskBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid Json")
		return
	}

	task := entity.Task{
		Name:   &body.Name,
		Status: &body.Status,
	}

	data, err := rH.h.CreateTask(ctx, task)
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Create Task Failed")
		return
	}

	ctx.WithData(data).Response(http.StatusOK, "")
}

type updateTaskBody struct {
	Name   string `json:"name" binding:"required,gt=0,lte=255"`
	Status uint8  `json:"status" binding:"required,oneof=0 1"`
}

func (rH *HttpHandler) UpdateTaskHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	idStr := ctx.Param("taskID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid ID")
		return
	}

	var body updateTaskBody
	if err = ctx.ShouldBindJSON(&body); err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid Json")
		return
	}

	data, err := rH.h.UpdateTask(ctx, uint(id), entity.Task{
		Name:   &body.Name,
		Status: &body.Status,
	})
	if err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Update Task Failed")
		return
	}

	ctx.WithData(data).Response(http.StatusOK, "")
}

func (rH *HttpHandler) DeleteTaskHandler(c *gin.Context) {
	ctx := cGin.NewContext(c)

	idStr := ctx.Param("taskID")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.WithError(err).Response(http.StatusBadRequest, "Invalid ID")
		return
	}

	if err := rH.h.DeleteTask(ctx, uint(id)); err != nil {
		ctx.WithError(err).Response(http.StatusInternalServerError, "Delete Task Failed")
		return
	}

	ctx.Response(http.StatusOK, "")
}
