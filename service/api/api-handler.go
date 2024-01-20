package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	rt.router.PUT("/user/:username", rt.wrap(rt.setMyUserName))
	rt.router.GET("/user/:username/profile", rt.wrap(rt.getUserProfile))
	rt.router.GET("/user/:username/id/:userid", rt.wrap(rt.getUsername))
	rt.router.GET("/user/:username/stream", rt.wrap(rt.getUserFollowedPhotos))

	rt.router.POST("/user/:username/photo", rt.wrap(rt.uploadPhoto))
	rt.router.DELETE("/user/:username/photo/:photoid", rt.wrap(rt.deletePhoto))
	rt.router.GET("/user/:username/photo", rt.wrap(rt.getUserPhotos))

	rt.router.POST("/user/:username/ban", rt.wrap(rt.banUser))
	rt.router.DELETE("/user/:username/ban/:banid", rt.wrap(rt.unbanUser))
	rt.router.GET("/user/:username/ban", rt.wrap(rt.getBans))

	rt.router.POST("/user/:username/follow", rt.wrap(rt.followUser))
	rt.router.DELETE("/user/:username/follow/:followid", rt.wrap(rt.unfollowUser))
	rt.router.GET("/user/:username/follow", rt.wrap(rt.getFollower))

	rt.router.POST("/user/:username/photo/:photoid/like", rt.wrap(rt.likePhoto))
	rt.router.DELETE("/user/:username/photo/:photoid/like/:likeid", rt.wrap(rt.unlikePhoto))
	rt.router.GET("/user/:username/photo/:photoid/like", rt.wrap(rt.getLike))

	rt.router.POST("/user/:username/photo/:photoid/comment", rt.wrap(rt.commentPhoto))
	rt.router.DELETE("/user/:username/photo/:photoid/comment/:commentid", rt.wrap(rt.uncommentPhoto))
	rt.router.GET("/user/:username/photo/:photoid/comment", rt.wrap(rt.getComment))

	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
