package example

import (
	"bytes"
	"context"
	"go/types"
	. "io"
)

func NamedContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithCancel(ctx)
}

func NamedBytes(ctx context.Context) Writer {
	var b bytes.Buffer
	return &b
}

func NamedStdFile() types.Importer {
	return nil
}
