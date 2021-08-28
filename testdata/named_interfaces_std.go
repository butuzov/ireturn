package example

import (
	"context"
	"go/types"
)

func NamedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

func NamedStdFile() types.Importer {
	return nil
}
