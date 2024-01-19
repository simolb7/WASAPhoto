<script>

export default {
    data: function () {
        return {
            errormsg: null,
            //username: localStorage.getItem('username'),
            token: localStorage.getItem('token'),
            newUsername: "",
            profile: {
                RequestId: 0,
                UserId: 0,
                Username: "",
                NumberFollowers: 0,
                NumberFollowed: 0,
                photoCount: 0,
            },
            user: {
                id: 0,
                username: "",
            },
            photoList: {
                requestUser: 0,
                identifier: 0,
                photos: [
                    {
                        id: 0,
                        userId: 0,
                        file: "",
                        date: "",
                        LikeNumber: 0,
                        CommentNumber: 0,
                        comment: "",
                        isUnLikeButton: false,
                        likeId: 0,
                        
                    }
                ],
            },
            photoComments: {
                requestIdentifier: 0,
                photoIdentifier: 0,
                identifier: 0,
                comments: [
                    {
                        id: 0,
                        userId: 0,
                        photoId: 0,
                        photoOwnerID: 0,
                        content: "",
                        ownerUsername: "",
                    }
                ],
            },
            follow: {
                followId: 0,
                followedId: 0,
                userId: 0,
                banStatus: 0,
            },
            ban: {
                banId: 0,
                bannedId: 0,
                userId: 0,
            },
        }
    },
    methods: {

        async refresh() {
            await this.userProfile();
            await this.userPhotos();
            
        },
        async userProfile() {
            try {
                let response = await this.$axios.get("/user/" + this.$route.params.username + "/profile", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.profile = response.data
                let followid = await this.getFollow();
                if (followid !== 0){
                    this.follow.followId = followid
                }
                let banid = await this.getBan();
                if (banid !== 0){
                    this.ban.banId = banid
                }
               
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async userPhotos() {
            try {
                let response = await this.$axios.get("/user/" + this.$route.params.username + "/photo", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.photoList = response.data
                this.photoList.photos.sort((a, b) => b.id - a.id);
                
                for (let i = 0; i < this.photoList.photos.length; i++) {
                    this.photoList.photos[i].file = 'data:image/jpg;base64,' + this.photoList.photos[i].file;
                
                    let likestatus = await this.getLikeStatus(localStorage.getItem('username'), this.photoList.photos[i].id);
                   
                    if (likestatus.hasLike) {
                        // Se esiste già un like, disabilita il pulsante "Like" e abilita il pulsante "Unlike"
                        this.photoList.photos[i].isUnlikeButton = true;
                        this.photoList.photos[i].likeId = likestatus.likeid;

                    } else {
                        this.photoList.photos[i].isUnlikeButton = false;     
                    }
                }
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    if (this.photoList && this.photoList.photos && this.photoList.photos.length > 0) {
                        // Ci sono foto, non dovresti vedere il messaggio di errore
                    } else {
                        this.errormsg = this.$route.params.username + " hasn't posted any photos yet. Go to the home and upload one!";
                        this.detailedmsg = null;
                    }
                }
            }
        },
        formatDateTime(dateTime) {
            //timeZoneName: 'short' 
            const options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit'};
            const formattedDateTime = new Date(dateTime).toLocaleString('en-US', options);
            return formattedDateTime;
        },
        async sendComment(photoid, comment, ownerid) {
            if (comment === null || comment === undefined) {
                this.errormsg = "Emtpy comment field."
            } else if(comment.length > 200) {
                this.errormsg = "Comment is too long."
            } else {
                try {
                    let response = await this.$axios.post("/user/" + localStorage.getItem('username') + "/photo/" + photoid + "/comment", { content: comment, photoOwnerID: ownerid }, {
                        headers: {
                            Authorization: "Bearer " + localStorage.getItem("token")
                        }
                    })
                    this.clear = response.data
                    this.refresh()
                } catch (e) {
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                        this.detailedmsg = null;
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                        this.detailedmsg = e.toString();
                    } else {
                        this.errormsg = e.toString();
                        this.detailedmsg = null;
                    }
                }
            }
        },
        async getComments(photoid) {
            try {
                let response = await this.$axios.get("/user/" + this.$route.params.username  + "/photo/" + photoid + "/comment", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.photoComments = response.data;
                if (this.photoComments && this.photoComments.comments && this.photoComments.comments.length > 0){
                    for (let i = 0; i < this.photoComments.comments.length; i++) {
                        this.photoComments.comments[i].ownerUsername = await this.getusername(this.photoComments.comments[i].userId);
                    }
                }
                const modal = new bootstrap.Modal(document.getElementById('commentPopup'));
                modal.show();
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async deleteComment(commentid, photoid, useridcommento) {
			try {
				let response = await this.$axios.delete("/user/" + localStorage.getItem('username') + "/photo/" + photoid + "/comment/" + commentid, {
					headers: {
						Authorization: "Bearer " + useridcommento
					}
				})
				location.reload();
			} catch(e) {
				if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
					this.detailedmsg = null;
				} else if (e.response && e.response.status === 500) {
					this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
					this.detailedmsg = e.toString();
				} else {
					this.errormsg = e.toString();
					this.detailedmsg = null;
				}
			}
		},
        canDeleteComment(comment) {
        // Supponiamo che tu abbia informazioni sull'utente autenticato e sulla proprietà della foto
            //console.log('comment:', comment);
            const isAuthenticatedUser = localStorage.getItem("username") === comment.ownerUsername;
            const isPhotoOwner = comment.photoOwnerID === parseInt(localStorage.getItem("token"));
            /*console.log('comment ID:', comment.id);
            console.log('id user auth:', localStorage.getItem("token"));
            console.log('id owner:', comment.photoOwnerID);
            console.log('isAuthenticatedUser:', isAuthenticatedUser);
            console.log('isPhotoOwner:', isPhotoOwner);
            */// Ritorna true solo se l'utente è autenticato e ha il permesso di eliminare il commento
            return isAuthenticatedUser || isPhotoOwner;
        },
        async getusername(userid) {
            try {
                let response = await this.$axios.get("/user/"+ this.username + "/id/" + userid , {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("token")}`
                }
                });
                return response.data.username;
            } catch (error) {
                if (error.response && error.response.status === 404) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
					this.detailedmsg = null;
                } 
            }
            console.error("Errore durante il recupero dello stato del like:", error);
        },
        async getLikeStatus(username, photoid) {
            try {
                let response = await this.$axios.get("/user/" + username + "/photo/" + photoid + "/like", {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("token")}`
                }
                });
                if (response.data !== null) {
                    return {
                        hasLike: true,
                        likeid: response.data.likeId
                    };
                } else {
                    return {
                        hasLike: false,
                        likeid: 0
                    };
        }
            } catch (error) {
                if (error.response && error.response.status === 404) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
					this.detailedmsg = null;
                } 
            }
            console.error("Errore durante il recupero dello stato del like:", error);
        },
        async likePhoto(ownerid,photoid) {
            try {
                console.log(ownerid, photoid)
                let response = await this.$axios.post("/user/" + localStorage.getItem('username') + "/photo/" + photoid + "/like", {userId: ownerid}, {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.clear = response.data
                this.refresh()
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async unlikePhoto(photoid, likeid) {
            try {
                    let response = await this.$axios.delete("/user/" + localStorage.getItem('username') + "/photo/" + photoid + "/like/" + likeid, {
                        headers: {
                            Authorization: "Bearer " + localStorage.getItem("token")
                        }
                    })
                    this.clear = response.data
                    this.refresh()
                } catch (e) {
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                        this.detailedmsg = null;
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                        this.detailedmsg = e.toString();
                    } else {
                        this.errormsg = e.toString();
                        this.detailedmsg = null;
                    }
                }
        },
        async followUser() {
            try {
                let response = await this.$axios.post("/user/" + this.$route.params.username + "/follow", {}, {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.clear = response.data
                this.refresh()
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async unfollowUser(followid) {
            try {
                let response = await this.$axios.delete("/user/" + this.$route.params.username + "/follow/" + followid , {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.clear = response.data
                this.follow.followId = 0;
                this.refresh()
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }

        },
        async getFollow() {
            try {
                let response = await this.$axios.get("/user/" + this.$route.params.username + "/follow", {
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("token")}`
                }
                });
                if (response.data !== null) {
                    this.follow = response.data
                    return response.data.followId
                    }
                else {
                    return 0
                }
            } catch (error) {
                if (error.response && error.response.status === 404) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
					this.detailedmsg = null;
                } 
            }
            console.error("Errore durante il recupero dello stato del follow:", error);
        },
        async banUser() {
            try {
                let response = await this.$axios.post("/user/" + this.$route.params.username + "/ban", {}, {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.clear = response.data
                this.refresh()
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async unbanUser() {
            try {
                let response = await this.$axios.delete("/user/" + this.$route.params.username + "/ban/" + this.ban.banId, {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.clear = response.data
                this.ban.banId = 0
                this.refresh()
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = e.toString();
                    this.detailedmsg = null;
                }
            }
        },
        async getBan() {
            try {
                let response = await this.$axios.get("/user/" + this.$route.params.username + "/ban", { //simone
                headers: {
                    Authorization: `Bearer ${localStorage.getItem("token")}` //Marta
                }
                });
                if (response.data !== null) {
                    this.ban = response.data
                    return response.data.banId
                    }
                else {
                    return 0
                }
            } catch (error) {
                if (error.response && error.response.status === 404) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
					this.detailedmsg = null;
                } 
            }
            console.error("Errore durante il recupero dello stato del ban:", error);
        },
        

	},

    mounted() {
        this.userProfile()
        this.userPhotos()
    }
}
</script>

<template>

    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-5 border-bottom">
        <h1 class="h2">Profile of <b>{{ profile.username }}</b>  </h1>
        <div class="p-4 text-black">
            <div v-if = "ban.banId === 0" class="d-flex justify-content-end text-center py-1">
                <div>
                    <p class="mb-1 h5">{{ profile.NumberFollowers }}</p>
                    <p class="small text-muted mb-0">Followers</p>
                </div>
                <div class="px-3">
                    <p class="mb-1 h5">{{ profile.NumberFollowed }}</p>
                    <p class="small text-muted mb-0">Followings</p>
                </div>
                <div class="mb-3 mx-5">
                    <p class="mb-1 h5">{{ profile.photoCount }}</p>
                    <p class="small text-muted mb-0">Photos</p>
                </div>
                <div>
                    <button v-if = "follow.followId === 0" class="btn btn-primary" type="button" @click="followUser">Follow</button>
                    <button v-if = "follow.followId !== 0" class="btn btn-danger" type="button" @click="unfollowUser(follow.followId)">Unfollow</button>
                </div>
                <div class="mb-3 mx-5">
                    <button class="btn btn-danger" type="button" @click="banUser">Ban</button>
                </div>
            </div>
        </div>
    </div>

    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
        v-if="ban.banId !== 0 && ban.userId !== token">
        
        <div class="alert alert-danger text-center mx-auto" role="alert">
            <h4 class="alert-heading">You have banned @{{ profile.username }}</h4>
            <p>Remove ban to interact.</p>
            <div >
                    <button class="btn btn-primary" type="button" @click="unbanUser">Remove ban</button>
            </div>
            <hr>
            <p class="mb-0"></p>
        </div>
    </div>
    <div class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom"
        v-if="ban.banId !== 0 && ban.userId === token">
        <div class="alert alert-danger text-center mx-auto" role="alert">
            <h4 class="alert-heading">Impossible perform this action...</h4>
            <p>User @{{ profile.username }} has banned you.</p>
            <p>You are not allowed to interact.</p>
            <hr>
            <p class="mb-0"></p>
        </div>
    </div>

    
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div class="row" v-if="photoList && photoList.photos && photoList.photos.length > 0 && ban.banId === 0 ">
        <div class="col-md-4" v-for="photo in photoList.photos" :key="photo.id">
            <div class="card mb-4 shadow-sm">
                <img class="card-img-top" :src=photo.file alt="Card image cap">
                <div class="card-body">
                    <div class="button-container">
                        <button type="button" class="btn btn-outline-primary expanded-button">@{{ profile.username }}</button>
                    </div>
                   
                    <div
                        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
                    </div>
                    <div class="d-flex justify-content-between align-items-center">
                        <p class="card-text">Likes : {{ photo.likeNumber }}</p>
                    </div>
                    <div class="d-flex justify-content-between align-items-center">
                        <p class="card-text">Comments : {{ photo.commentNumber }}</p>
                    </div>
                    <p class="card-text">Photo uploaded on {{ formatDateTime(photo.date) }}</p>

                    <div class="input-group mb-3">
                        <input type="text" id="comment" v-model="photo.comment" class="form-control" placeholder="Comment!"
                          aria-describedby="basic-addon2">
                        <div class="input-group-append">
                            <button class="btn btn-primary" type="button"
                                @click="sendComment(photo.id, photo.comment, photo.userId)">Send</button>
                        </div>
                    </div>

                    <div class="d-flex justify-content-between align-items-center">
                        <div class="btn-group">
                            <!--<button type="button" class="btn btn-dark" @click="openLog(username, photo.id)">Comments</button>--> 
                            <button v-if="photo.isUnlikeButton" type="button" class="btn btn-danger" @click="unlikePhoto(photo.id, photo.likeId)">
                                Unlike
                                </button>
                            <button v-if="!photo.isUnlikeButton" type="button" class="btn btn-primary" @click="likePhoto(photo.userId, photo.id)">
                                Like
                            </button>
                            
                        </div>
                        <button type="button" class="btn btn-success" @click="getComments(photo.id)">View all comments</button>
         
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div id="commentPopup" class="modal fade">followUser
        <div class="modal-dialog">
            <div class="modal-content">
                <!-- Contenuto della finestra di popup -->
                <div class="modal-header">
                    <h5 class="modal-title">All comments</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <ul v-if="photoComments && photoComments.comments && photoComments.comments.length > 0" style = "list-style-type: none;">
                        <li v-for="comment in photoComments.comments" :key="comment.id" class="comment-item ">
                            <div  class="d-flex justify-content-between align-items-center">
                                <div class="comment-container"> 
                                    <strong>{{ comment.ownerUsername }}</strong> 
                                    <p class = "comment-content">{{ comment.content }}</p>
                                </div>
                                <div class = "ml-auto"  v-if="canDeleteComment(comment)" >
                                    <button class="btn btn-danger mr-2" @click="deleteComment(comment.id, comment.photoId, comment.userId)">Delete</button>
                                </div>
                                
                            </div>
                        </li>
                    </ul>
                    <p v-else> No comments available.</p>
                </div>
            </div>
        </div>
    </div>


   
</template>

<style>

     .button-container {
        text-align: center; /* Centra il pulsante orizzontalmente */
    }
    .expanded-button {
        display: block;
        width: 100%;
        margin-top: 10px; /* Aggiungi margine in alto se necessario */
    }
    .comment-item {
        margin-bottom: 10px; /* Aggiungi margine inferiore tra i commenti */
    }
    .comment-content {
        overflow-wrap: break-word;
        word-wrap: break-word;
        word-break: break-word;
    }
    .comment-container {
        /* Stile per il contenitore del commento */
        margin-right: 8px; /* O qualsiasi valore desiderato */
    }
</style>