package api

import (
	"encoding/json"
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

	//utente da followare
	//username := ps.ByName("username")

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userFollowedID, ok := requestData["username"]
	if !ok {
		http.Error(w, "Il campo 'userFollowedId' è obbligatorio", http.StatusBadRequest)
		return
	}

	dbuser, err := rt.db.GetUserId(userFollowedID)
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
	//username := ps.ByName("username")

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userFollowed, ok := requestData["username"]
	if !ok {
		http.Error(w, "Il campo 'userFollowedId' è obbligatorio", http.StatusBadRequest)
		return
	}

	dbuser, err := rt.db.GetUserId(userFollowed)
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
		//registro un messaggio di errore, aggiungo un campo "id" al log per indicare l'ID del commrnto
		http.Error(w, fmt.Sprintf("Errore durante l'eliminazione del follow con ID %d", id), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getFollower(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var follow Follow
	token := getToken(r.Header.Get("Authorization"))

	var requestData map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userFollowed, ok := requestData["username"]
	if !ok {
		http.Error(w, "Il campo 'username' è obbligatorio", http.StatusBadRequest)
		return
	}
	user.Username = userFollowed
	dbuser, err := rt.db.CheckUserByUsername(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	dbfollow, err := rt.db.GetFollower(user.ToDatabase(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	follow.FollowFromDatabase(dbfollow)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(follow)
}
