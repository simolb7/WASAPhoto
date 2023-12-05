package api

import (
	"encoding/json"
	"net/http"

	//"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/simolb7/WASAPhoto/service/api/reqcontext"
	"github.com/simolb7/WASAPhoto/service/database"
)

// This function provide to log the user or create a new user by an username
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User

	//decodifica la richiesta http
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	dbuser, err := rt.db.CreateUser(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// this function provide to change user's username
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	username := ps.ByName("username")

	if username == "" {
		http.Error(w, "L'username non può essere vuoto", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token := getToken(r.Header.Get("Authorization"))

	user.Id = token
	dbuser, err := rt.db.SetUsername(user.ToDatabase(), username)
	if err != nil {
		//gestisco gli errori specifici
		if err == database.ErrUsernameAlreadyExists {
			http.Error(w, "Username già in uso", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	user.FromDatabase(dbuser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// This function return user profile
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User
	var requestUser User
	var profile Profile
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
	followersCount, err := rt.db.GetFollowersCount(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	followedCount, err := rt.db.GetFollowedCount(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	photoCount, err := rt.db.GetPhotosCount(user.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	profile.RequestId = token
	profile.UserId = user.Id
	profile.Username = user.Username
	profile.NumberFollowers = followersCount
	profile.NumberFollowed = followedCount
	profile.PhotoCount = photoCount

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(profile)
}

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// create user struct
	var user User
	//	create database stream struct
	var photoList database.Stream
	// Ottieni l'ID dell'utente dal token di autorizzazione
	token := getToken(r.Header.Get("Authorization"))
	// Ottieni il nome utente dalla URL
	username := ps.ByName("username")

	// Creare la struttura User con le informazioni ottenute
	user.Id = token
	user.Username = username

	// Verifica l'utente
	dbuser, err := rt.db.CheckUser(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.FromDatabase(dbuser)

	// Ottieni lo stream di foto dell'utente dal database
	stream, err := rt.db.GetStream(user.ToDatabase())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Creare la struttura Stream con le informazioni ottenute
	photoList.Id = token
	photoList.Photos = stream

	// Imposta l'intestazione e restituisci lo stream in formato JSON con uno stato HTTP 200 (OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(photoList)
}

/*
func (rt *_router) getUserUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Ottieni l'ID dalla richiesta
	id := ps.ByName("id")
	var user User

	if id == "" {
		http.Error(w, "L'ID non può essere vuoto", http.StatusBadRequest)
		return
	}

	// Converti l'ID in un tipo appropriato (ad esempio, int)
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Formato ID non valido", http.StatusBadRequest)
		return
	}

	// Ottieni l'username dal database usando l'ID
	dbuser, err := rt.db.GetUserById(idInt)
	if err != nil {
		// Gestisci gli errori specifici
		if err == database.ErrUserDoesNotExist {
			http.Error(w, "Utente non trovato", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	user.FromDatabase(dbuser)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}
*/
