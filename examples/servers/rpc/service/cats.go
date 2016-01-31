package service

import (
	"net/http"

	"github.com/xhroot/gizmo/examples/nyt"
	"github.com/xhroot/gizmo/server"
	"golang.org/x/net/context"
)

func (s *RPCService) GetCats(ctx context.Context, r *CatsRequest) (*CatsResponse, error) {
	var (
		err error
		res []*nyt.SemanticConceptArticle
	)
	defer server.MonitorRPCRequest()(ctx, "GetCats", err)

	res, err = s.client.SemanticConceptSearch("des", "cats")
	if err != nil {
		return nil, err
	}

	return &CatsResponse{res}, nil
}

func (s *RPCService) GetCatsJSON(r *http.Request) (int, interface{}, error) {
	res, err := s.GetCats(context.Background(), &CatsRequest{})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}
	return http.StatusOK, res.Results, nil
}
