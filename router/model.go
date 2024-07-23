package router

import "github.com/Lynicis/inzibat/config"

const ErrorTypeCasting = "type casting error"

type RouteChannel struct {
	IndexOfRoute int
	Route        config.Route
}
