package api

import (
	"api-test-coverage/internal/utils"
	"context"
)

func GetTime(ctx context.Context) *ResponseV1 {
	timeNow := utils.GetTime(ctx)
	return &ResponseV1{
		Status: true,
		Data: ResponseV1Data{
			Time: timeNow,
		},
	}
}
