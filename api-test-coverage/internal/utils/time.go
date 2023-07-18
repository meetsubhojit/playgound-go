package utils

import (
	"context"
	"time"
)

func GetTime(ctx context.Context) string {
	return time.Now().String()
}
