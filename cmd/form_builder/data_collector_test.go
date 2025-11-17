package form_builder

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"inzibat/config"
)

func TestHeadersCollector(t *testing.T) {
	t.Run("happy path - GetEmptyValue returns empty header", func(t *testing.T) {
		collector := &HeadersCollector{}

		emptyValue := collector.GetEmptyValue()

		assert.NotNil(t, emptyValue)
		headers, ok := emptyValue.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("happy path - GetSourceTitle returns correct title", func(t *testing.T) {
		collector := &HeadersCollector{}

		title := collector.GetSourceTitle()

		assert.Equal(t, "Header Source", title)
	})

	t.Run("happy path - GetFileFormConfig returns correct config", func(t *testing.T) {
		collector := &HeadersCollector{}

		config := collector.GetFileFormConfig()

		assert.Equal(t, "filepath", config.Key)
		assert.Equal(t, "Header JSON File Path", config.Title)
		assert.Equal(t, "/path/to/headers.json", config.Placeholder)
	})

	t.Run("happy path - CollectFromFile loads headers", func(t *testing.T) {
		collector := &HeadersCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "headers.json")
		headersData := map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer token",
		}
		data, err := json.Marshal(headersData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
		assert.Equal(t, "Bearer token", headers.Get("Authorization"))
	})

	t.Run("happy path - CollectFromForm calls CollectHeadersFromForm", func(t *testing.T) {
		collector := &HeadersCollector{}

		result, err := collector.CollectFromForm()

		if err == nil {
			assert.NotNil(t, result)
			headers, ok := result.(http.Header)
			assert.True(t, ok)
			assert.NotNil(t, headers)
		} else {
			assert.Error(t, err)
		}
	})
}

func TestBodyCollector(t *testing.T) {
	t.Run("happy path - GetEmptyValue returns nil HttpBody", func(t *testing.T) {
		collector := &BodyCollector{}

		emptyValue := collector.GetEmptyValue()

		body, ok := emptyValue.(config.HttpBody)
		assert.True(t, ok)
		if body != nil {
			assert.Equal(t, 0, len(body))
		}
	})

	t.Run("happy path - GetSourceTitle returns correct title", func(t *testing.T) {
		collector := &BodyCollector{}

		title := collector.GetSourceTitle()

		assert.Equal(t, "Body Source", title)
	})

	t.Run("happy path - GetFileFormConfig returns correct config", func(t *testing.T) {
		collector := &BodyCollector{}

		config := collector.GetFileFormConfig()

		assert.Equal(t, "filepath", config.Key)
		assert.Equal(t, "Body JSON File Path", config.Title)
		assert.Equal(t, "/path/to/body.json", config.Placeholder)
	})

	t.Run("happy path - CollectFromFile loads body", func(t *testing.T) {
		collector := &BodyCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.json")
		bodyData := config.HttpBody{
			"message": "success",
			"code":    float64(200),
		}
		data, err := json.Marshal(bodyData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Equal(t, "success", body["message"])
		assert.Equal(t, float64(200), body["code"])
	})

	t.Run("happy path - CollectFromForm calls CollectBodyFromForm", func(t *testing.T) {
		collector := &BodyCollector{}

		result, err := collector.CollectFromForm()

		if err == nil {
			_, ok := result.(config.HttpBody)
			assert.True(t, ok)
		} else {
			assert.Error(t, err)
		}
	})
}

func TestBodyStringCollector(t *testing.T) {
	t.Run("happy path - GetEmptyValue returns empty string", func(t *testing.T) {
		collector := &BodyStringCollector{}

		emptyValue := collector.GetEmptyValue()

		assert.Equal(t, "", emptyValue)
	})

	t.Run("happy path - GetSourceTitle returns correct title", func(t *testing.T) {
		collector := &BodyStringCollector{}

		title := collector.GetSourceTitle()

		assert.Equal(t, "BodyString Source", title)
	})

	t.Run("happy path - GetFileFormConfig returns correct config", func(t *testing.T) {
		collector := &BodyStringCollector{}

		config := collector.GetFileFormConfig()

		assert.Equal(t, FilePathKey, config.Key)
		assert.Equal(t, "BodyString File Path", config.Title)
		assert.Equal(t, "/path/to/body.txt", config.Placeholder)
	})

	t.Run("happy path - CollectFromFile loads body string", func(t *testing.T) {
		collector := &BodyStringCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.txt")
		expectedContent := `{"message": "success"}`
		err := os.WriteFile(filePath, []byte(expectedContent), 0644)
		require.NoError(t, err)

		result, err := collector.CollectFromFile(filePath)

		assert.NoError(t, err)
		assert.Equal(t, expectedContent, result)
	})

	t.Run("happy path - CollectFromForm calls CollectBodyStringFromForm", func(t *testing.T) {
		collector := &BodyStringCollector{}

		result, err := collector.CollectFromForm()

		if err == nil {
			_, ok := result.(string)
			assert.True(t, ok)
		} else {
			assert.Error(t, err)
		}
	})
}

func TestCollectData(t *testing.T) {
	t.Run("happy path - skip source returns empty value", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceSkip).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return NewMockFormRunner(ctrl)
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("happy path - file source collects from file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "headers.json")
		headersData := map[string]string{
			"Content-Type": "application/json",
		}
		data, err := json.Marshal(headersData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(nil).Times(1)
		filePathMock.EXPECT().GetString("filepath").Return(filePath).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		assert.NotNil(t, result)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, "application/json", headers.Get("Content-Type"))
	})

	t.Run("happy path - form source collects from form", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}

		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceForm).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return NewMockFormRunner(ctrl)
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		if err == nil {
			assert.NotNil(t, result)
		} else {
			assert.Error(t, err)
		}
	})

	t.Run("happy path - body collector with file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &BodyCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.json")
		bodyData := config.HttpBody{
			"test": "value",
		}
		data, err := json.Marshal(bodyData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(nil).Times(1)
		filePathMock.EXPECT().GetString("filepath").Return(filePath).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Equal(t, "value", body["test"])
	})

	t.Run("happy path - body string collector with file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &BodyStringCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.txt")
		expectedContent := "test content"
		err := os.WriteFile(filePath, []byte(expectedContent), 0644)
		require.NoError(t, err)

		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(nil).Times(1)
		filePathMock.EXPECT().GetString(FilePathKey).Return(filePath).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		assert.Equal(t, expectedContent, result)
	})

	t.Run("error path - source form run fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(errors.New("form run failed")).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return NewMockFormRunner(ctrl)
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("error path - file path form run fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(errors.New("file path form failed")).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("error path - collect from file fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(nil).Times(1)
		filePathMock.EXPECT().GetString("filepath").Return("/nonexistent/file.json").Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCollectHeaders(t *testing.T) {
	t.Run("happy path - returns headers collector interface", func(t *testing.T) {
		collector := &HeadersCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})

	t.Run("happy path - collects headers with skip source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceSkip).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return NewMockFormRunner(ctrl)
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, 0, len(headers))
	})

	t.Run("happy path - collects headers with file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &HeadersCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "headers.json")
		headersData := map[string]string{
			"X-Test": "value",
		}
		data, err := json.Marshal(headersData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(nil).Times(1)
		filePathMock.EXPECT().GetString("filepath").Return(filePath).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		headers, ok := result.(http.Header)
		assert.True(t, ok)
		assert.Equal(t, "value", headers.Get("X-Test"))
	})
}

func TestCollectBody(t *testing.T) {
	t.Run("happy path - returns body collector interface", func(t *testing.T) {
		collector := &BodyCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})

	t.Run("happy path - collects body with skip source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &BodyCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceSkip).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return NewMockFormRunner(ctrl)
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Nil(t, body)
	})

	t.Run("happy path - collects body with file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &BodyCollector{}
		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "body.json")
		bodyData := config.HttpBody{
			"message": "test",
		}
		data, err := json.Marshal(bodyData)
		require.NoError(t, err)
		err = os.WriteFile(filePath, data, 0644)
		require.NoError(t, err)

		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceFile).Times(1)

		filePathMock := NewMockFormRunner(ctrl)
		filePathMock.EXPECT().Run().Return(nil).Times(1)
		filePathMock.EXPECT().GetString("filepath").Return(filePath).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return filePathMock
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		assert.NotNil(t, result)
		body, ok := result.(config.HttpBody)
		assert.True(t, ok)
		assert.Equal(t, "test", body["message"])
	})
}

func TestCollectBodyString(t *testing.T) {
	t.Run("happy path - returns body string collector interface", func(t *testing.T) {
		collector := &BodyStringCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})

	t.Run("happy path - collects body string with skip source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		collector := &BodyStringCollector{}
		sourceMock := NewMockFormRunner(ctrl)
		sourceMock.EXPECT().Run().Return(nil).Times(1)
		sourceMock.EXPECT().GetString("source").Return(SourceSkip).Times(1)

		getSourceForm := func() FormRunner {
			return sourceMock
		}
		getFilePathForm := func() FormRunner {
			return NewMockFormRunner(ctrl)
		}

		result, err := collectDataWithRunners(collector, getSourceForm, getFilePathForm)

		require.NoError(t, err)
		bodyString, ok := result.(string)
		assert.True(t, ok)
		assert.Equal(t, "", bodyString)
	})
}

func TestCollectData_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectData with headers collector and skip source", func(t *testing.T) {
		collector := &HeadersCollector{}
		assert.Implements(t, (*DataCollector)(nil), collector)
	})

	t.Run("happy path - CollectData with body collector", func(t *testing.T) {
		collector := &BodyCollector{}

		assert.Implements(t, (*DataCollector)(nil), collector)
	})

	t.Run("happy path - CollectData with body string collector", func(t *testing.T) {
		collector := &BodyStringCollector{}

		assert.Implements(t, (*DataCollector)(nil), collector)
	})

	t.Run("happy path - CollectData function exists and creates forms", func(t *testing.T) {
		collector := &HeadersCollector{}

		assert.Implements(t, (*DataCollector)(nil), collector)
	})
}

func TestCollectHeaders_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectHeaders function exists and returns correct type", func(t *testing.T) {
		collector := &HeadersCollector{}

		assert.Implements(t, (*DataCollector)(nil), collector)
		assert.Equal(t, "Header Source", collector.GetSourceTitle())
	})
}

func TestCollectBody_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectBody function exists and returns correct type", func(t *testing.T) {
		collector := &BodyCollector{}

		assert.Implements(t, (*DataCollector)(nil), collector)
		assert.Equal(t, "Body Source", collector.GetSourceTitle())
	})
}

func TestCollectBodyString_PublicFunction(t *testing.T) {
	t.Run("happy path - CollectBodyString function exists and returns correct type", func(t *testing.T) {
		collector := &BodyStringCollector{}

		assert.Implements(t, (*DataCollector)(nil), collector)
		assert.Equal(t, "BodyString Source", collector.GetSourceTitle())
	})
}
