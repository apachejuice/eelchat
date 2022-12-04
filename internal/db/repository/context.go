package repository

import (
	"context"
	"time"
)

func doCtx(closure func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()
	return closure(ctx)
}
