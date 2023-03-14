package processor

import "github.com/go-yaaf/yaaf-code-gen/model"

// Processor interface
type Processor interface {
	Process(metaModel *model.MetaModel) error
	Name() string
}
