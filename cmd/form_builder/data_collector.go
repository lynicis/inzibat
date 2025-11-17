package form_builder

import (
	"net/http"

	"inzibat/config"
)

func CollectHeaders() (http.Header, error) {
	sourceForm := BuildSourceSelectionForm("Header Source", SourceKey)
	if err := sourceForm.Run(); err != nil {
		return nil, err
	}
	source := sourceForm.GetString(SourceKey)

	if source == SourceSkip {
		return make(http.Header), nil
	}

	if source == SourceFile {
		filePathForm := BuildFilePathForm(FilePathFormConfig{
			Key:         "filepath",
			Title:       "Header JSON File Path",
			Placeholder: "/path/to/headers.json",
		})
		if err := filePathForm.Run(); err != nil {
			return nil, err
		}
		filePath := filePathForm.GetString("filepath")
		return config.LoadHeadersFromFile(filePath)
	}

	return CollectHeadersFromForm()
}

func CollectBody() (config.HttpBody, error) {
	sourceForm := BuildSourceSelectionForm("Body Source", SourceKey)
	if err := sourceForm.Run(); err != nil {
		return nil, err
	}
	source := sourceForm.GetString(SourceKey)

	if source == SourceSkip {
		return nil, nil
	}

	if source == SourceFile {
		filePathForm := BuildFilePathForm(FilePathFormConfig{
			Key:         "filepath",
			Title:       "Body JSON File Path",
			Placeholder: "/path/to/body.json",
		})
		if err := filePathForm.Run(); err != nil {
			return nil, err
		}
		filePath := filePathForm.GetString("filepath")
		return config.LoadBodyFromFile(filePath)
	}

	return CollectBodyFromForm()
}

func CollectBodyString() (string, error) {
	sourceForm := BuildSourceSelectionForm("BodyString Source", SourceKey)
	if err := sourceForm.Run(); err != nil {
		return "", err
	}
	source := sourceForm.GetString(SourceKey)

	if source == SourceSkip {
		return "", nil
	}

	if source == SourceFile {
		filePathForm := BuildFilePathForm(FilePathFormConfig{
			Key:         FilePathKey,
			Title:       "BodyString File Path",
			Placeholder: "/path/to/body.txt",
		})
		if err := filePathForm.Run(); err != nil {
			return "", err
		}
		filePath := filePathForm.GetString(FilePathKey)
		return config.LoadBodyStringFromFile(filePath)
	}

	return CollectBodyStringFromForm()
}
