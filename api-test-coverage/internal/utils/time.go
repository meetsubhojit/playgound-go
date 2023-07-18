package utils

import (
	"context"
	"time"
)

func GetTime(ctx context.Context) string {
	return time.Now().String()
}

func NotUsedFunction() int {
	return 0
}
