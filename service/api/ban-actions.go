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

	dbuser, err := rt.db.GetUserId(ps.ByName("username"))
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
	err = rt.db.RemoveFollows(token, user.Id)
	if err != nil && !errors.Is(err, database.ErrFollowDoesNotExist) {
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

	dbuser, err := rt.db.GetUserId(ps.ByName("username"))
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

	user.Username = ps.ByName("username")
	dbuser, err := rt.db.CheckUserByUsername(user.ToDatabase())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.FromDatabase(dbuser)
	dban, err := rt.db.GetBan(user.ToDatabase(), token)
	if errors.Is(err, database.ErrBanDoesNotExist) {
		// Il ban non esiste, restituisci null

		_, err1 := w.Write([]byte("null"))
		if err1 != nil {
			// Gestisci l'errore in qualche modo
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	ban.BanFromDatabase(dban)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ban)
}
