package handler

import (
	"net/http"
	"encoding/json"
	"log"

	"github.com/TechBowl-japan/go-stations/model"
)

type PanicHandler struct{}

func NewPanicHandler() *PanicHandler {
	return &PanicHandler{}
}

func (h *PanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	panic("Panic!")

	// 以下は実行されないことを確認
	res := &model.HealthzResponse{Message: "Panicking"}
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		log.Println(err)
	}
}
