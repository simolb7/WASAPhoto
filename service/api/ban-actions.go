package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolb7/WASAPhoto/service/api/reqcontext"
	"github.com/simolb7/WASAPhoto/service/database"
)

func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ban Ban
	var user User

	token := getToken(r.Header.Get("Authorization"))

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userban, ok := requestData["username"]
	if !ok {
		http.Error(w, "Il campo 'username' è obbligatorio", http.StatusBadRequest)
		return
	}

	dbuser, err := rt.db.GetUserId(userban)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	ban.BannedId = user.Id
	ban.UserId = token
	dbban, err := rt.db.InsertBan(ban.BanToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ban.BanFromDatabase(dbban)

	err = rt.db.RemoveComments(token, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = rt.db.RemoveLikes(token, user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(ban)
}

func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var ban Ban
	var user User
	token := getToken(r.Header.Get("Authorization"))

	id, err := strconv.ParseUint(ps.ByName("banid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var requestData map[string]string
	err1 := json.NewDecoder(r.Body).Decode(&requestData)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	userban, ok := requestData["username"]
	if !ok {
		http.Error(w, "Il campo 'username' è obbligatorio", http.StatusBadRequest)
		return
	}

	dbuser, err := rt.db.GetUserId(userban)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)
	ban.BanId = id
	ban.UserId = token
	ban.BannedId = user.Id
	err = rt.db.RemoveBan(ban.BanToDatabase())
	if errors.Is(err, database.ErrBanDoesNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Errore durante l'eliminazione del ban con ID %d", id), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getBans(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var ban Ban
	token := getToken(r.Header.Get("Authorization"))

	var requestData map[string]string
	err1 := json.NewDecoder(r.Body).Decode(&requestData)
	if err1 != nil {
		http.Error(w, err1.Error(), http.StatusBadRequest)
		return
	}

	userban, ok := requestData["username"]
	if !ok {
		http.Error(w, "Il campo 'username' è obbligatorio", http.StatusBadRequest)
		return
	}

	user.Username = userban
	dbuser, err := rt.db.CheckUserByUsername(user.ToDatabase())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.FromDatabase(dbuser)
	dban, err := rt.db.GetBan(user.ToDatabase(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	ban.BanFromDatabase(dban)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ban)
}