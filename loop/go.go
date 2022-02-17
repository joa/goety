package loop

import (
	"context"
)

// Go - Call f in a loop until the context has been cancelled.
//
// It doesn't matter if f will panic; it will be continuously called
// until the context is done.
func Go(ctx context.Context, f func()) {
	GoErr(ctx, nil, f)
}

// GoErr - Call f in a loop until the context has been cancelled.
//
// It doesn't matter if f will panic; it will be continuously called
// until the context is done.
//
// All observed errors will be sent to the errs channel. It's optional
// and no error will be propagated if it is nil.
func GoErr(ctx context.Context, errs chan interface{}, f func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if errs != nil {
					select {
					case <-ctx.Done():
						return
					case errs <- err:
					}
				}

				GoErr(ctx, errs, f)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				f()
			}
		}
	}()
}
