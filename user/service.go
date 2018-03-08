package user

import "context"

type Service interface {
	// GET /user/<uid>
	Show(ctx context.Context, uid string)
}
