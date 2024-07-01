package router

import (
	"net/http"
	"strconv"

	"github.com/ars0915/gogolook-exercise/util/cGin"
)

func resourceCheck(rH HttpHandler) cGin.HandlerFunc {
	return func(ctx *cGin.Context) {
		taskIDStr := ctx.Param("taskID")

		if len(taskIDStr) > 0 {
			taskID, err := strconv.Atoi(taskIDStr)
			if err != nil {
				ctx.WithError(err).Response(http.StatusBadRequest, "Invalid ID")
				return
			}

			if _, err = rH.Usecase().GetTask(ctx, uint(taskID)); err != nil {
				ctx.WithError(err).Response(http.StatusInternalServerError, "Check Task failed")
				return
			}
		}

		ctx.Next()
	}
}
