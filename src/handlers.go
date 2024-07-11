package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

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

// ReturnStatus provides addition information about response
type ReturnStatus struct {
	Error   string      `json:"error,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

func onPage(ectx echo.Context) error {
	var address string
	if err := echo.QueryParamsBinder(ectx).
		String("url", &address).
		BindError(); err != nil {
		return ectx.JSON(http.StatusBadRequest, ReturnStatus{Error: fmt.Sprintf("url query parameter is missing or invalid: %s", err.Error())})
	}

	// query 'url'
	tr := &http.Transport{
		MaxIdleConns:       1,
		IdleConnTimeout:    10 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(address)
	if err != nil {
		return ectx.JSON(http.StatusBadRequest, ReturnStatus{Error: fmt.Sprintf("Failed to read page at provided address '%s': %s", address, err.Error())})
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ectx.JSON(http.StatusInternalServerError, ReturnStatus{Error: err.Error()})
	}
	return ectx.JSON(resp.StatusCode, ReturnStatus{Payload: string(bodyBytes)})
}
