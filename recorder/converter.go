package recorder

import (
	"net/http"

	"github.com/goccy/go-json"

	"github.com/lynicis/inzibat/config"
)

// ConvertToInzibatConfig converts a recorded session into an inzibat mock configuration.
// Duplicate (method, path) pairs are deduplicated — the last recording wins.
func ConvertToInzibatConfig(session RecordedSession, serverPort int) *config.Cfg {
	if serverPort <= 0 {
		serverPort = 8080
	}

	routeMap := map[string]*config.Route{}
	routeOrder := []string{}

	for _, entry := range session.Entries {
		key := entry.Request.Method + " " + entry.Request.Path

		route := &config.Route{
			Method:       entry.Request.Method,
			Path:         entry.Request.Path,
			FakeResponse: buildFakeResponse(entry.Response),
		}

		if _, exists := routeMap[key]; !exists {
			routeOrder = append(routeOrder, key)
		}

		routeMap[key] = route
	}

	routes := make([]config.Route, 0, len(routeOrder))
	for _, key := range routeOrder {
		routes = append(routes, *routeMap[key])
	}

	return &config.Cfg{
		ServerPort: serverPort,
		Routes:     routes,
	}
}

func buildFakeResponse(resp RecordedResponse) *config.FakeResponse {
	fakeResponse := &config.FakeResponse{
		StatusCode: resp.StatusCode,
		Headers:    convertHeaders(resp.Headers),
	}

	if len(resp.Body) == 0 {
		return fakeResponse
	}

	// Try to unmarshal as a JSON object for Body (map[string]any)
	var bodyMap config.HttpBody
	if err := json.Unmarshal(resp.Body, &bodyMap); err == nil && bodyMap != nil {
		fakeResponse.Body = bodyMap
		return fakeResponse
	}

	// Fall back to BodyString
	var bodyStr string
	if err := json.Unmarshal(resp.Body, &bodyStr); err == nil {
		fakeResponse.BodyString = bodyStr
		return fakeResponse
	}

	// Raw fallback: store the JSON literal as a string
	fakeResponse.BodyString = string(resp.Body)

	return fakeResponse
}

func convertHeaders(headers map[string][]string) http.Header {
	if len(headers) == 0 {
		return nil
	}

	result := make(http.Header, len(headers))
	for key, values := range headers {
		result[key] = values
	}

	return result
}
