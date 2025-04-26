package i

import "context"

type ProductRepository interface {
	Repository
	UpdateStock(ctx context.Context, productID uint, quantity int) error
}
