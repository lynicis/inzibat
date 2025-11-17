package form_builder

import (
	"net/http"

	"inzibat/config"
)

// collectHeadersInternal is the testable internal implementation that accepts form runners
func collectHeadersInternal(sourceFormRunner, filePathFormRunner FormRunner) (http.Header, error) {
	if err := sourceFormRunner.Run(); err != nil {
		return nil, err
	}
	source := sourceFormRunner.GetString(SourceKey)

	if source == SourceSkip {
		return make(http.Header), nil
	}

	if source == SourceFile {
		if err := filePathFormRunner.Run(); err != nil {
			return nil, err
		}
		filePath := filePathFormRunner.GetString("filepath")
		return config.LoadHeadersFromFile(filePath)
	}

	return CollectHeadersFromForm()
}

func CollectHeaders() (http.Header, error) {
	sourceForm := BuildSourceSelectionForm("Header Source", SourceKey)
	filePathForm := BuildFilePathForm(FilePathFormConfig{
		Key:         "filepath",
		Title:       "Header JSON File Path",
		Placeholder: "/path/to/headers.json",
	})

	sourceFormRunner := &huhFormRunner{form: sourceForm}
	filePathFormRunner := &huhFormRunner{form: filePathForm}

	return collectHeadersInternal(sourceFormRunner, filePathFormRunner)
}

// collectBodyInternal is the testable internal implementation that accepts form runners
func collectBodyInternal(sourceFormRunner, filePathFormRunner FormRunner) (config.HttpBody, error) {
	if err := sourceFormRunner.Run(); err != nil {
		return nil, err
	}
	source := sourceFormRunner.GetString(SourceKey)

	if source == SourceSkip {
		return nil, nil
	}

	if source == SourceFile {
		if err := filePathFormRunner.Run(); err != nil {
			return nil, err
		}
		filePath := filePathFormRunner.GetString("filepath")
		return config.LoadBodyFromFile(filePath)
	}

	return CollectBodyFromForm()
}

func CollectBody() (config.HttpBody, error) {
	sourceForm := BuildSourceSelectionForm("Body Source", SourceKey)
	filePathForm := BuildFilePathForm(FilePathFormConfig{
		Key:         "filepath",
		Title:       "Body JSON File Path",
		Placeholder: "/path/to/body.json",
	})

	sourceFormRunner := &huhFormRunner{form: sourceForm}
	filePathFormRunner := &huhFormRunner{form: filePathForm}

	return collectBodyInternal(sourceFormRunner, filePathFormRunner)
}

// collectBodyStringInternal is the testable internal implementation that accepts form runners
func collectBodyStringInternal(sourceFormRunner, filePathFormRunner FormRunner) (string, error) {
	if err := sourceFormRunner.Run(); err != nil {
		return "", err
	}
	source := sourceFormRunner.GetString(SourceKey)

	if source == SourceSkip {
		return "", nil
	}

	if source == SourceFile {
		if err := filePathFormRunner.Run(); err != nil {
			return "", err
		}
		filePath := filePathFormRunner.GetString(FilePathKey)
		return config.LoadBodyStringFromFile(filePath)
	}

	return CollectBodyStringFromForm()
}

func CollectBodyString() (string, error) {
	sourceForm := BuildSourceSelectionForm("BodyString Source", SourceKey)
	filePathForm := BuildFilePathForm(FilePathFormConfig{
		Key:         FilePathKey,
		Title:       "BodyString File Path",
		Placeholder: "/path/to/body.txt",
	})

	sourceFormRunner := &huhFormRunner{form: sourceForm}
	filePathFormRunner := &huhFormRunner{form: filePathForm}

	return collectBodyStringInternal(sourceFormRunner, filePathFormRunner)
}
