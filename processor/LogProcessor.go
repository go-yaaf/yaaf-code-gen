package processor

import (
	"encoding/json"
	"fmt"
	"github.com/go-yaaf/yaaf-code-gen/model"
	"os"
)

// LogProcessor - Log processor write model to a file for inspection
type LogProcessor struct {
	outputFolder string
}

// NewLogProcessor - Factory method
func NewLogProcessor(outputFolder string) Processor {
	return &LogProcessor{outputFolder: outputFolder}
}

// Name returns the processor name
func (p *LogProcessor) Name() string {
	return "Log Processor"
}

// Process starts the processor
func (p *LogProcessor) Process(metaModel *model.MetaModel) error {

	// First, ensure output directory
	if err := os.MkdirAll(p.outputFolder, os.ModePerm); err != nil {
		return fmt.Errorf("error creating folder: %s: %s", p.outputFolder, err)
	}

	// For debugging purpose, print the model
	if bytes, err := json.MarshalIndent(metaModel, "", "  "); err != nil {
		return fmt.Errorf("error converting model to json: %s", err)
	} else {

		// print
		fmt.Println(string(bytes))

		// save
		filePath := fmt.Sprintf("%s/model.json", p.outputFolder)
		if f, fer := os.Create(filePath); fer != nil {
			return fmt.Errorf("error creating file: %s: %s", filePath, fer)
		} else {
			if _, wrr := f.Write(bytes); err != nil {
				return fmt.Errorf("error writing to file: %s: %s", filePath, wrr)
			} else {
				return f.Close()
			}

		}
	}
}
