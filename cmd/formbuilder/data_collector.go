package formbuilder

import (
	"net/http"

	"github.com/lynicis/inzibat/config"
)

func collectFromSource[T any](
	sourceFormRunner,
	filePathFormRunner FormRunner,
	loadFile func(string) (T, error),
	collectFromForm func() (T, error),
	skipValue T,
) (T, error) {
	var zero T

	if err := sourceFormRunner.Run(); err != nil {
		return zero, err
	}

	switch sourceFormRunner.GetString(SourceKey) {
	case SourceSkip:
		return skipValue, nil
	case SourceFile:
		if err := filePathFormRunner.Run(); err != nil {
			return zero, err
		}

		return loadFile(filePathFormRunner.GetString(FilePathKey))
	default:
		return collectFromForm()
	}
}

func collectHeadersInternal(
	sourceFormRunner,
	filePathFormRunner FormRunner,
) (http.Header, error) {
	return collectFromSource(
		sourceFormRunner,
		filePathFormRunner,
		config.LoadHeadersFromFile,
		CollectHeadersFromForm,
		make(http.Header),
	)
}

func CollectHeaders() (http.Header, error) {
	sourceForm := BuildSourceSelectionForm("Header Source", SourceKey)
	filePathForm := BuildFilePathForm(FilePathFormConfig{
		Key:         FilePathKey,
		Title:       "Header JSON File Path",
		Placeholder: "/path/to/headers.json",
	})

	sourceFormRunner := &HuhFormRunner{Form: sourceForm}
	filePathFormRunner := &HuhFormRunner{Form: filePathForm}

	return collectHeadersInternal(sourceFormRunner, filePathFormRunner)
}

func collectBodyInternal(
	sourceFormRunner,
	filePathFormRunner FormRunner,
) (config.HttpBody, error) {
	return collectFromSource(
		sourceFormRunner,
		filePathFormRunner,
		config.LoadBodyFromFile,
		CollectBodyFromForm,
		config.HttpBody(nil),
	)
}

func CollectBody() (config.HttpBody, error) {
	sourceForm := BuildSourceSelectionForm("Body Source", SourceKey)
	filePathForm := BuildFilePathForm(FilePathFormConfig{
		Key:         FilePathKey,
		Title:       "Body JSON File Path",
		Placeholder: "/path/to/body.json",
	})

	sourceFormRunner := &HuhFormRunner{Form: sourceForm}
	filePathFormRunner := &HuhFormRunner{Form: filePathForm}

	return collectBodyInternal(sourceFormRunner, filePathFormRunner)
}

func collectBodyStringInternal(
	sourceFormRunner,
	filePathFormRunner FormRunner,
) (string, error) {
	return collectFromSource(
		sourceFormRunner,
		filePathFormRunner,
		config.LoadBodyStringFromFile,
		CollectBodyStringFromForm,
		"",
	)
}

func CollectBodyString() (string, error) {
	sourceForm := BuildSourceSelectionForm("BodyString Source", SourceKey)
	filePathForm := BuildFilePathForm(FilePathFormConfig{
		Key:         FilePathKey,
		Title:       "BodyString File Path",
		Placeholder: "/path/to/body.txt",
	})

	sourceFormRunner := &HuhFormRunner{Form: sourceForm}
	filePathFormRunner := &HuhFormRunner{Form: filePathForm}

	return collectBodyStringInternal(sourceFormRunner, filePathFormRunner)
}
