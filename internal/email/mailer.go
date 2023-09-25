package email

import "context"

type Mailer interface {
	Send(ctx context.Context, email *Email) error
}
