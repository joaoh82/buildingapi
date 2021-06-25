package presenters

import (
	"net/http"

	"github.com/joaoh82/buildingapi/utils"
)

func (p *presenters) JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	err := utils.WriteJson(w, v)
	if err != nil {
		p.Error(w, r, err)
	}
}
