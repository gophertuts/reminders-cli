package controllers

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gophertuts/reminders-cli/server/models"
)

// ctx param fetches param from context
func ctxParam(ctx context.Context, key string) urlParam {
	ps, ok := ctx.Value(ctxKey(paramsKey)).(map[string]urlParam)
	if !ok {
		return urlParam{}
	}
	return ps[key]
}

// parseIDParam parses id url param
func parseIDParam(ctx context.Context) (int, error) {
	id, err := strconv.Atoi(ctxParam(ctx, idParamName).value)
	if err != nil {
		return 0, models.DataValidationError{Message: "invalid id provided"}
	}
	return id, nil
}

// parseIDParam parses ids url param
func parseIDsParam(ctx context.Context) ([]int, error) {
	idsSlice := strings.Split(ctxParam(ctx, idsParamName).value, ",")
	var res []int
	var invalid []int
	for _, id := range idsSlice {
		n, err := strconv.Atoi(id)
		if err != nil {
			invalid = append(invalid, n)
		}
		res = append(res, n)
	}
	if len(invalid) > 0 {
		err := models.DataValidationError{
			Message: fmt.Sprintf("invalid ids provided: %v", invalid),
		}
		return []int{}, err
	}
	return res, nil
}
