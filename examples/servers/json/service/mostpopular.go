package service

import (
	"net/http"

	"github.com/xhroot/gizmo/web"
	"github.com/gorilla/mux"
)

func (s *JSONService) GetMostPopular(r *http.Request) (int, interface{}, error) {
	resourceType := mux.Vars(r)["resourceType"]
	section := mux.Vars(r)["section"]
	timeframe := web.GetUInt64Var(r, "timeframe")
	res, err := s.client.GetMostPopular(resourceType, section, uint(timeframe))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res, nil
}
