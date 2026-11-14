package handler

import (
	"fmt"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	httpPkg "github.com/lynicis/inzibat/client/http"
	"github.com/lynicis/inzibat/config"
)

type ClientHandler struct {
	Client                  *httpPkg.Client
	RouteConfig             *[]config.Route
	CircuitBreakerStore     *CircuitBreakerStore
	CircuitBreakerRouteKeys map[int]string
}

func (clientRoute *ClientHandler) CreateHandler(routeIndex int) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		requestTo := (*clientRoute.RouteConfig)[routeIndex].RequestTo
		routeKey, hasCircuitBreaker := clientRoute.CircuitBreakerRouteKeys[routeIndex]

		isAllowed, err := clientRoute.allowRequest(hasCircuitBreaker, routeKey)
		if err != nil {
			return err
		}
		if !isAllowed {
			return ctx.
				Status(fiber.StatusServiceUnavailable).
				SendString("circuit breaker is open")
		}

		parsedUrl, err := requestTo.GetParsedUrl()
		if err != nil {
			return fmt.Errorf("failed to parse request URL: %w", err)
		}

		var bodyBytes []byte
		if len(requestTo.Body) > 0 {
			bodyBytes, err = json.Marshal(&requestTo.Body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
		}

		response, err := clientRoute.Client.Do(
			requestTo.Method,
			parsedUrl.String(),
			requestTo.Headers,
			bodyBytes,
		)
		if err != nil {
			if recordErr := clientRoute.recordFailure(hasCircuitBreaker, routeKey); recordErr != nil {
				return recordErr
			}
			if requestTo.InErrorReturn500 {
				ctx.Status(fiber.StatusInternalServerError)
				return ctx.Send(nil)
			}

			return ctx.
				Status(fiber.StatusInternalServerError).
				SendString(err.Error())
		}

		if recordErr := clientRoute.recordSuccess(hasCircuitBreaker, routeKey); recordErr != nil {
			return recordErr
		}

		return ctx.
			Status(response.Status).
			Send(response.Body)
	}
}

func (clientRoute *ClientHandler) allowRequest(hasCircuitBreaker bool, routeKey string) (bool, error) {
	if !hasCircuitBreaker {
		return true, nil
	}

	allowed, err := clientRoute.CircuitBreakerStore.Allow(routeKey)
	if err != nil {
		return false, fmt.Errorf("failed to evaluate circuit breaker: %w", err)
	}

	return allowed, nil
}

func (clientRoute *ClientHandler) recordFailure(hasCircuitBreaker bool, routeKey string) error {
	if !hasCircuitBreaker {
		return nil
	}

	err := clientRoute.CircuitBreakerStore.OnFailure(routeKey)
	if err != nil {
		return fmt.Errorf("failed to update circuit breaker failure: %w", err)
	}

	return nil
}

func (clientRoute *ClientHandler) recordSuccess(hasCircuitBreaker bool, routeKey string) error {
	if !hasCircuitBreaker {
		return nil
	}

	err := clientRoute.CircuitBreakerStore.OnSuccess(routeKey)
	if err != nil {
		return fmt.Errorf("failed to update circuit breaker success: %w", err)
	}

	return nil
}
