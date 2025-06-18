package document

import "context"

// go get github.com/baliance/gooxml/document - word
// go get github.com/ledongthuc/pdf -pdf

type DocumentReader interface {
	ReadDocument(ctx context.Context, path string) (string, error)
}
