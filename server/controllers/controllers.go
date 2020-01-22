package controllers

import (
	"context"
	"log"
	"strconv"
	"strings"
)

// ctx param fetches param from context
func ctxParam(ctx context.Context, key string) urlParam {
	ps, ok := ctx.Value(ctxKey("ps")).(map[string]urlParam)
	if !ok {
		return urlParam{}
	}
	return ps[key]
}

func parseIDParam(ctx context.Context) int {
	id, err := strconv.Atoi(ctxParam(ctx, idParamName).value)
	if err != nil {
		log.Fatalf("could not convert id to number: %v", err)
	}
	return id
}

func parseIDsParam(ctx context.Context) []int {
	idsSlice := strings.Split(ctxParam(ctx, idsParamName).value, ",")
	var res []int
	for _, id := range idsSlice {
		n, err := strconv.Atoi(id)
		if err != nil {
			log.Fatalf("could not convert id to number: %v", err)
		}
		res = append(res, n)
	}
	return res
}
