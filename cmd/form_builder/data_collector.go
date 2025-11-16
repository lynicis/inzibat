package form_builder

import (
	"net/http"

	"inzibat/config"
)

type DataCollector interface {
	CollectFromFile(filePath string) (interface{}, error)
	CollectFromForm() (interface{}, error)
	GetEmptyValue() interface{}
	GetSourceTitle() string
	GetFileFormConfig() FilePathFormConfig
}

type HeadersCollector struct{}

func (c *HeadersCollector) CollectFromFile(filePath string) (interface{}, error) {
	return config.LoadHeadersFromFile(filePath)
}

func (c *HeadersCollector) CollectFromForm() (interface{}, error) {
	return CollectHeadersFromForm()
}

func (c *HeadersCollector) GetEmptyValue() interface{} {
	return make(http.Header)
}

func (c *HeadersCollector) GetSourceTitle() string {
	return "Header Source"
}

func (c *HeadersCollector) GetFileFormConfig() FilePathFormConfig {
	return FilePathFormConfig{
		Key:         "filepath",
		Title:       "Header JSON File Path",
		Placeholder: "/path/to/headers.json",
	}
}

type BodyCollector struct{}

func (c *BodyCollector) CollectFromFile(filePath string) (interface{}, error) {
	return config.LoadBodyFromFile(filePath)
}

func (c *BodyCollector) CollectFromForm() (interface{}, error) {
	return CollectBodyFromForm()
}

func (c *BodyCollector) GetEmptyValue() interface{} {
	return config.HttpBody(nil)
}

func (c *BodyCollector) GetSourceTitle() string {
	return "Body Source"
}

func (c *BodyCollector) GetFileFormConfig() FilePathFormConfig {
	return FilePathFormConfig{
		Key:         "filepath",
		Title:       "Body JSON File Path",
		Placeholder: "/path/to/body.json",
	}
}

type BodyStringCollector struct{}

func (c *BodyStringCollector) CollectFromFile(filePath string) (interface{}, error) {
	return config.LoadBodyStringFromFile(filePath)
}

func (c *BodyStringCollector) CollectFromForm() (interface{}, error) {
	return CollectBodyStringFromForm()
}

func (c *BodyStringCollector) GetEmptyValue() interface{} {
	return ""
}

func (c *BodyStringCollector) GetSourceTitle() string {
	return "BodyString Source"
}

func (c *BodyStringCollector) GetFileFormConfig() FilePathFormConfig {
	return FilePathFormConfig{
		Key:         FilePathKey,
		Title:       "BodyString File Path",
		Placeholder: "/path/to/body.txt",
	}
}

func collectDataWithRunners(
	collector DataCollector,
	getSourceForm func() FormRunner,
	getFilePathForm func() FormRunner,
) (interface{}, error) {
	sourceForm := getSourceForm()
	if err := sourceForm.Run(); err != nil {
		return nil, err
	}
	source := sourceForm.GetString("source")

	if source == SourceSkip {
		return collector.GetEmptyValue(), nil
	}

	if source == SourceFile {
		filePathForm := getFilePathForm()
		if err := filePathForm.Run(); err != nil {
			return nil, err
		}
		filePath := filePathForm.GetString(collector.GetFileFormConfig().Key)
		return collector.CollectFromFile(filePath)
	}

	return collector.CollectFromForm()
}

func CollectData(collector DataCollector) (interface{}, error) {
	return collectDataWithRunners(
		collector,
		func() FormRunner {
			return &huhFormRunner{
				form: BuildSourceSelectionForm(collector.GetSourceTitle(), SourceKey),
			}
		},
		func() FormRunner {
			return &huhFormRunner{
				form: BuildFilePathForm(collector.GetFileFormConfig()),
			}
		},
	)
}

func CollectHeaders() (http.Header, error) {
	collector := &HeadersCollector{}
	result, err := CollectData(collector)
	if err != nil {
		return nil, err
	}
	return result.(http.Header), nil
}

func CollectBody() (config.HttpBody, error) {
	collector := &BodyCollector{}
	result, err := CollectData(collector)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, nil
	}
	return result.(config.HttpBody), nil
}

func CollectBodyString() (string, error) {
	collector := &BodyStringCollector{}
	result, err := CollectData(collector)
	if err != nil {
		return "", err
	}
	return result.(string), nil
}
