package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolb7/WASAPhoto/service/api/reqcontext"
	"github.com/simolb7/WASAPhoto/service/database"
)

// Create a like istance, the photOwnerID is taken from the body
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var like Like
	type PhotoOwnerStruct struct {
		UserId int `json:"userId"`
	}

	var requestBody PhotoOwnerStruct

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	photoid, err := strconv.ParseUint(ps.ByName("photoid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := getToken(r.Header.Get("Authorization"))

	like.UserId = token
	like.PhotoId = photoid
	like.PhotoOwner = uint64(requestBody.UserId)

	dblike, err := rt.db.InsertLike(like.LikeToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	like.LikeFromDatabase(dblike)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(like)
}

// Delete a like with a specific id from a specific photo, both taken from the path
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var like Like
	username := ps.ByName("username")
	dbuser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)
	photoid, err := strconv.ParseUint(ps.ByName("photoid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	likeid, err := strconv.ParseUint(ps.ByName("likeid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	like.LikeId = likeid
	token := getToken(r.Header.Get("Authorization"))
	like.UserId = token
	like.PhotoId = photoid
	like.PhotoOwner = user.Id
	err = rt.db.RemoveLike(like.LikeToDatabase())
	if errors.Is(err, database.ErrLikeDoesNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).WithField("id", likeid).Error("impossible delete this like")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Returns a like by a specific user, if there isn't a like, return "null"
func (rt *_router) getLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var user User
	var photo Photo
	var requestUser User
	token := getToken(r.Header.Get("Authorization"))
	requestUser.Id = token
	dbrequestuser, err := rt.db.CheckUserById(requestUser.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestUser.FromDatabase(dbrequestuser)
	username := ps.ByName("username")
	dbuser, err := rt.db.GetUserId(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.FromDatabase(dbuser)
	photoid, err := strconv.ParseUint(ps.ByName("photoid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	photo.Id = photoid

	dbphoto, err := rt.db.CheckPhoto(photo.PhotoToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	photo.PhotoFromDatabase(dbphoto)
	like, err := rt.db.GetLike(photo.Id, token)

	if errors.Is(err, database.ErrLikeDoesNotExist) {
		// Il like non esiste, restituisci null
		_, err1 := w.Write([]byte("null"))
		if err1 != nil {

			http.Error(w, "Failed to write response", http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Se il like esiste, restituisci le informazioni del like come JSON
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(like)
}
