package service

import (
	"encoding/json"
	"net/http"

	"github.com/xhroot/gizmo/examples/nyt"
)

func (s *JSONPubService) PublishCats(r *http.Request) (int, interface{}, error) {
	var catArticle nyt.SemanticConceptArticle
	err := json.NewDecoder(r.Body).Decode(&catArticle)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	err = s.pub.Publish(catArticle.Url, &catArticle)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	res := struct {
		Status string `json:"status"`
	}{
		"success!",
	}
	return http.StatusOK, res, nil
}
