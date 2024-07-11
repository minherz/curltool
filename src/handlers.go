package main

import (
	"github.com/labstack/echo/v4"
)

type RouteDefition struct {
	Name    string
	Path    string
	Handler echo.HandlerFunc
}

var RoutingMap = []RouteDefition{
	{
		Name:    "page-retriever",
		Path:    "/page",
		Handler: onPage,
	},
}

func onPage(ectx echo.Context) error {
	Logger.Info("page handler is called")
	return nil
}
