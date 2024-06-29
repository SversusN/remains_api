package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"golang.org/x/text/encoding/unicode"
	"io"
	"net/http"
	"remains_api/internal/domain"
	"remains_api/internal/repository"
	"remains_api/pkg/auth"
)

type Handlers struct {
	s  repository.Storage
	au *jwtauth.JWTAuth
}

func NewHandlers(s repository.Storage, au *jwtauth.JWTAuth) *Handlers {
	return &Handlers{s, au}
}

func (h *Handlers) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	token, _, _ := jwtauth.FromContext(r.Context())
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if token == nil {
		http.Error(w, "почистите кэш", http.StatusUnauthorized)
		return
	}

	userid, ok := token.Get("userID")
	if !ok {
		http.Error(w, "почистите куки", http.StatusUnauthorized)
		return
	}
	fmt.Println(userid)
	results, err := h.s.GetAll(userid.(string))
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

func (h *Handlers) LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	params := domain.LoginStruct{}
	body, err := io.ReadAll(r.Body)
	err = json.Unmarshal(body, &params)
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	userInfo, err := h.s.LoginUser(params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	userInfoString, _ := encoder.String(userInfo)
	token := auth.MakeToken(userInfoString)
	fmt.Println(userInfoString)
	_, err = w.Write([]byte(token))

}
