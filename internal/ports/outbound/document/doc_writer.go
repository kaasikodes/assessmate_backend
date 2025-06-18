package document

// go get github.com/baliance/gooxml/document - for word
// go get github.com/jung-kurt/gofpdf - for pdf
type DocumentGenerator interface {
	GenerateDocx(data any, path string) error
	GeneratePDF(data any, path string) error
}
