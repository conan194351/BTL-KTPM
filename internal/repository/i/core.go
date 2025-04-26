package i

import "context"

type Repository interface {
	GetByID(ctx context.Context, id uint) (interface{}, error)
	Create(ctx context.Context, user interface{}) error
}
