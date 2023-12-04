package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.PUT("/user/:username/setusername", rt.wrap(rt.setMyUserName))
	rt.router.GET("/user/:username/stream", rt.wrap(rt.getMyStream))
	rt.router.GET("/user/:username/profile", rt.wrap(rt.getUserProfile))
	//rt.router.GET("/session", rt.wrap(rt.GetUserId))

	rt.router.POST("/user/:username/photo/uploadPhoto", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/user/:username/photo/:photoid/deletePhoto", rt.wrap(rt.deletePhoto))
	rt.router.GET("/user/:username/photo/getUserPhotos", rt.wrap(rt.getUserPhotos))
	/*
		rt.router.PUT("/users/:username/ban/:banid", rt.wrap(rt.banUser))
		rt.router.DELETE("/users/:username/ban/:banid", rt.wrap(rt.unbanUser))
		rt.router.GET("/users/:username/ban", rt.wrap(rt.getBans))

		rt.router.PUT("/users/:username/follow/:followid", rt.wrap(rt.followUser))
		rt.router.DELETE("/users/:username/follow/:followid", rt.wrap(rt.unfollowUser))
		rt.router.GET("/users/:username/follow", rt.wrap(rt.getFollowers))
	*/
	rt.router.PUT("/users/:username/photo/:photoid/like/:likeid/likePhoto", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/users/:username/photo/:photoid/like/:likeid/unlikePhoto", rt.wrap(rt.unlikePhoto))
	rt.router.GET("/users/:username/photo/:photoid/like/getLike", rt.wrap(rt.getLike))
	/*
		rt.router.POST("/users/:username/photo/:photoid/comment", rt.wrap(rt.commentPhoto))
		rt.router.DELETE("/users/:username/photo/:photoid/comment/:commentid", rt.wrap(rt.uncommentPhoto))
		rt.router.GET("/users/:username/photo/:photoid/comment", rt.wrap(rt.getComments))
		rt.router.GET("/liveness", rt.liveness)
	*/
	return rt.router
}
