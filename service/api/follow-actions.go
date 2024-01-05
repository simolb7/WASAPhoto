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

func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var follow Follow
	var user User

	// utente da followare
	username := ps.ByName("username")

	dbuser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validazione del token di autorizzazione
	token := getToken(r.Header.Get("Authorization"))

	follow.UserFollowedId = user.Id
	follow.UserId = token

	dbfollow, err := rt.db.InsertFollow(follow.FollowToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Restituzione della risposta con l'istanza Follow creata
	follow.FollowFromDatabase(dbfollow)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(follow)
}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var follow Follow
	var user User

	// Validazione dell'utente da unfolloware
	username := ps.ByName("username")
	dbuser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	id, err := strconv.ParseUint(ps.ByName("followid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := getToken(r.Header.Get("Authorization"))
	println(id)
	follow.FollowId = id
	follow.UserFollowedId = user.Id
	follow.UserId = token

	// Rimozione dell'istanza Follow dal database
	err = rt.db.RemoveFollow(follow.FollowId, follow.UserId, follow.UserFollowedId)
	if err == database.ErrBanDoesNotExist {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		// registro un messaggio di errore, aggiungo un campo "id" al log per indicare l'ID del commrnto
		http.Error(w, fmt.Sprintf("Errore durante l'eliminazione del follow con ID %d", id), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getFollower(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var follow Follow
	token := getToken(r.Header.Get("Authorization"))

	user.Username = ps.ByName("username")
	dbuser, err := rt.db.CheckUserByUsername(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	dbfollow, err := rt.db.GetFollower(user.ToDatabase(), token)
	if errors.Is(err, database.ErrFollowDoesNotExist) {
		// Il follow non esiste, restituisci null
		w.Write([]byte("null"))
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	follow.FollowFromDatabase(dbfollow)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(follow)
}
