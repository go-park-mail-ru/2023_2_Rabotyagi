package usecases_

import "context"

type ICategoryStorage interface {
	GetFullCategories(ctx context.Context)
}
