package handlers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"remains_api/internal/domain"
	"remains_api/internal/repository"
)

type Handlers struct {
	s repository.Storage
}

func NewHandlers(s repository.Storage) *Handlers {
	return &Handlers{s}
}

func (h *Handlers) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	results, err := h.s.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(results) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	response, _ := json.Marshal(results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handlers) GetFilteredHandle(w http.ResponseWriter, r *http.Request) {
	params := domain.RemainRequest{}
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	results, err := h.s.GetFiltered(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(results) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	response, _ := json.Marshal(results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handlers) GetGroupHandler(w http.ResponseWriter, r *http.Request) {
	params := domain.RemainRequest{}
	key := chi.URLParam(r, "group")
	body, err := io.ReadAll(r.Body)
	json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	results, err := h.s.GetOnlyGroup(key, params)
	if len(results) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	response, _ := json.Marshal(results)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
