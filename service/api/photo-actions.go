package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/simolb7/WASAPhoto/service/api/reqcontext"
	"github.com/simolb7/WASAPhoto/service/database"
)

// Upload a photo, takes in input all the information and the populate a photo var.
// the file is taken from the body and the date is taken from the system
func (rt *_router) uploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var photo Photo

	token := getToken(r.Header.Get("Authorization"))

	user.Id = token
	user.Username = ps.ByName("username")

	// Verifica l'esistenza dell'utente nel database
	dbuser, err := rt.db.CheckUser(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	// leggo la foto e inizio a creare la struct photo
	photo.File, err = io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}
	romeLocation, err := time.LoadLocation("Europe/Rome")
	if err != nil {
		// Gestisci l'errore
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Impostaola data e l'ID dell'utente nella struct photo
	currentTime := time.Now().In(romeLocation)

	photo.DateTime = currentTime.Format(time.RFC3339)
	photo.UserId = user.Id

	dbphoto, err := rt.db.InsertPhoto(photo.PhotoToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	photo.PhotoFromDatabase(dbphoto)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(photo)
}

// Delete a photo with a specific id taken from the path
func (rt *_router) deletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	token := getToken(r.Header.Get("Authorization"))

	// prendo l'id della foto
	photoid, err := strconv.ParseUint(ps.ByName("photoid"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := ps.ByName("username")

	user.Username = username
	user.Id = token
	dbuser, err := rt.db.CheckUser(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	// elimino la foto
	err = rt.db.RemovePhoto(photoid)

	if errors.Is(err, database.ErrPhotoDoesNotExist) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if err != nil {
		// registro un messaggio di errore, aggiungo un campo "id" al log per indicare l'ID della foto
		http.Error(w, fmt.Sprintf("Errore durante l'eliminazione della foto con ID %d", photoid), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// return all the photo of an user (not the logged one) with the request user id and the id of the user searched
func (rt *_router) getUserPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var requestUser User
	var photoList database.Photos

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
	photos, err := rt.db.GetPhotos(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	photoList.RequestUser = requestUser.Id
	photoList.Identifier = token
	photoList.Photos = photos
	w.Header().Set("Content-Type", "image/*")
	_ = json.NewEncoder(w).Encode(photoList)
}

// returns a stream of photos taken from the followers of the logged user, sorted in descending date order
func (rt *_router) getUserFollowedPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var requestUser User
	var photoList database.Photos

	token := getToken(r.Header.Get("Authorization"))
	requestUser.Id = token

	dbrequestuser, err := rt.db.CheckUserById(requestUser.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestUser.FromDatabase(dbrequestuser)

	photos, err := rt.db.GetPhotosFollower(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	photoList.RequestUser = requestUser.Id
	photoList.Identifier = requestUser.Id
	photoList.Photos = photos
	w.Header().Set("Content-Type", "image/*")
	_ = json.NewEncoder(w).Encode(photoList)
}
