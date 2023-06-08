package handler

import (
	"context"
	"net/http"
	"encoding/json"
	"log"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: *todo}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todos, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return nil, err
	}
	return &model.ReadTODOResponse{TODOs: todos}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: *todo}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if r.Method == "GET" {
		var rreq model.ReadTODORequest
		q := r.URL.Query()

		var err error
		if q.Get("prev_id") == "" {
			rreq.PrevID = 0
		} else {
			rreq.PrevID, err = strconv.ParseInt(q.Get("prev_id"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		if q.Get("size") == "" {
			rreq.Size = 5
		} else {
			rreq.Size, err = strconv.ParseInt(q.Get("size"), 10, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		rres, err := h.Read(ctx, &rreq)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.NewEncoder(w).Encode(rres)
		if err != nil {
			log.Println(err)
		}
	} else if r.Method == "POST" {
		var c model.CreateTODORequest
		err := json.NewDecoder(r.Body).Decode(&c)
		if err != nil {
			log.Println(err)
			return
		} else if (c.Subject == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		cres, err := h.Create(ctx, &c)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.NewEncoder(w).Encode(cres)
		if err != nil {
			log.Println(err)
		}
	} else if r.Method == "PUT" {
		var u model.UpdateTODORequest
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			log.Println(err)
			return
		} else if (u.Subject == "") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		ures, err := h.Update(ctx, &u)
		if err != nil {
			switch err.(type) {
			case *model.ErrNotFound:
				w.WriteHeader(http.StatusNotFound)
			default:
				log.Println(err)
			}
			return
		}
		err = json.NewEncoder(w).Encode(ures)
		if err != nil {
			log.Println(err)
		}
	}
}