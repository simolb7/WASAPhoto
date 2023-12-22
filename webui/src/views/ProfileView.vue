<script>
//import LogModal from "../components/Logmodal.vue";

export default {
    //components: { LogModal },
    data: function () {
        return {
            errormsg: null,
            username: localStorage.getItem('username'),
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
                        //ownerUsername: "",
                        username: "",
                        content: "",
                    }
                ],
            },
        }
    },
    methods: {

        async refresh() {
            
			this.loading = true;
			this.errormsg = null;
			try {
				let response = this.$router.go({ path: '/user/' + this.username + '/profile' })
				this.guesses = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
            await this.userProfile();
        },
        async userProfile() {
            try {
                let response = await this.$axios.get("/user/" + this.username + "/profile", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.profile = response.data
                //console.log('Profile data:', this.profile); // Aggiunto il console.log
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
        async doLogout() {
			localStorage.removeItem("token")
			localStorage.removeItem("username")
			this.$router.push({ path: '/' })
		},
        async userPhotos() {
            try {
                let response = await this.$axios.get("/user/" + this.username + "/photo", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.photoList = response.data
                this.photoList.photos.sort((a, b) => b.id - a.id);
                for (let i = 0; i < this.photoList.photos.length; i++) {
                    this.photoList.photos[i].file = 'data:image/jpg;base64,' + this.photoList.photos[i].file
                }
            } catch (e) {
                if (e.response && e.response.status === 400) {
                    this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                    this.detailedmsg = null;
                } else if (e.response && e.response.status === 500) {
                    this.errormsg = "An internal error occurred. We will be notified. Please try again later.";
                    this.detailedmsg = e.toString();
                } else {
                    this.errormsg = "You haven't posted any photos yet. Go to the home and upload one!";
                    this.detailedmsg = null;
                }
            }
        },
        async changeName() {
            if (this.newUsername == "") {
                this.errormsg = "Emtpy username field."
            } else {
                try {
                    let response = await this.$axios.put("/user/" + this.username , { username: this.newUsername }, {
                        headers: {
                            Authorization: "Bearer " + localStorage.getItem("token")
                        }
                    })
                    this.user = response.data
                    localStorage.setItem("username", this.user.username);
                    this.profile.username = this.user.username;
                    this.username = this.user.username;
                    this.newUsername = "";
                    this.$router.push({ path: '/user/' + this.user.username + '/profile' })
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
        async deletePhoto(photoid) {
            try {
                let response = await this.$axios.delete("/user/" + this.username + "/photo/" + photoid, {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.refresh();
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
        formatDateTime(dateTime) {
            //timeZoneName: 'short' 
            const options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit'};
            const formattedDateTime = new Date(dateTime).toLocaleString('en-US', options);
            return formattedDateTime;
        },
        async sendComment(username, photoid, comment, ownerid) {
            if (comment === "") {
                this.errormsg = "Emtpy comment field."
            } else {
                try {
                    let response = await this.$axios.post("/user/" + username + "/photo/" + photoid + "/comment", { content: comment, photoOwnerID: ownerid }, {
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
        async getComments(username, photoid) {
            try {
                let response = await this.$axios.get("/user/" + username + "/photo/" + photoid + "/comment", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.photoComments = response.data;
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
        async deleteComment(commentid, photoid, username, useridcommento) {
			try {
				let response = await this.$axios.delete("/user/" + username + "/photo/" + photoid + "/comment/" + commentid, {
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
            const isAuthenticatedUser = localStorage.getItem("username") === comment.username;
            const isPhotoOwner = comment.photoOwnerID === parseInt(localStorage.getItem("token"));
            /*console.log('comment ID:', comment.id);
            console.log('id user auth:', localStorage.getItem("token"));
            console.log('id owner:', comment.photoOwnerID);
            console.log('isAuthenticatedUser:', isAuthenticatedUser);
            console.log('isPhotoOwner:', isPhotoOwner);
            */// Ritorna true solo se l'utente è autenticato e ha il permesso di eliminare il commento
            return isAuthenticatedUser || isPhotoOwner;
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
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
        <h1 class="h2">Profile of <b>{{ profile.username }}</b>  </h1>
        <div class="p-4 text-black">
            <div class="d-flex justify-content-end text-center py-1">
                <div>
                    <p class="mb-1 h5">{{ profile.followersCount }}</p>
                    <p class="small text-muted mb-0">Followers</p>
                </div>
                <div class="px-3">
                    <p class="mb-1 h5">{{ profile.followingCount }}</p>
                    <p class="small text-muted mb-0">Followings</p>
                </div>
                <div class="mb-3 mx-5">
                    <p class="mb-1 h5">{{ profile.photoCount }}</p>
                    <p class="small text-muted mb-0">Photos</p>
                </div>
                <div>
                    <button class="btn btn-danger" type="button" @click="doLogout">Logout</button>
                </div>
            </div>
        </div>
    </div>
    
   
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

    <div class="input-group mb-3">
        <input type="text" id="newUsername" v-model="newUsername" class="form-control"
            placeholder="Insert a new username for your profile" aria-label="Recipient's username"
            aria-describedby="basic-addon2">
        <div class="input-group-append">
            <button class="btn btn-success" type="button" @click="changeName">Change username</button>
        </div>
    </div>

    <div
        class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
    </div>

   <!--. <LogModal id="logviewer" :log="photoComments" :token="token"></LogModal>..-->

    <div class="row" v-if="photoList && photoList.photos && photoList.photos.length > 0">
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
                        <p class="card-text">Likes : {{ photo.likesCount }}</p>
                    </div>
                    <div class="d-flex justify-content-between align-items-center">
                        <p class="card-text">Comments : {{ photo.commentsCount }}</p>
                    </div>
                    <p class="card-text">Photo uploaded on {{ formatDateTime(photo.date) }}</p>

                    <div class="input-group mb-3">
                        <input type="text" id="comment" v-model="photo.comment" class="form-control" placeholder="Comment!"
                          aria-describedby="basic-addon2">
                        <div class="input-group-append">
                            <button class="btn btn-primary" type="button"
                                @click="sendComment(username, photo.id, photo.comment, photo.userId)">Send</button>
                        </div>
                    </div>

                    <div class="d-flex justify-content-between align-items-center">
                        <div class="btn-group">
                            <!--<button type="button" class="btn btn-dark" @click="openLog(username, photo.id)">Comments</button>--> 
                            <button type="button" class="btn btn-primary" @click="likePhoto(username, photo.id)">Like</button>
                            <button type="button" v-if="photo.likeStatus == true" class="btn btn-danger"
                                @click="deleteLike(username, photo.id)">Unlike</button>
                            
                            <button type="button" class="btn btn-success" @click="getComments(username, photo.id)">View all comments</button>
                        </div>

                            <button type="button" class="btn btn-danger ml-auto"
                                @click="deletePhoto(photo.id)">Delete Photo</button>
                        
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div id="commentPopup" class="modal fade">
        <div class="modal-dialog">
            <div class="modal-content">
                <!-- Contenuto della finestra di popup -->
                <div class="modal-header">
                    <h5 class="modal-title">Comments</h5>
                    <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
                </div>
                <div class="modal-body">
                    <ul v-if="photoComments && photoComments.comments && photoComments.comments.length > 0">
                        <li v-for="comment in photoComments.comments" :key="comment.id" >
                            <div  class="d-flex justify-content-between">
                                <div>
                                    <strong>{{ comment.username }}</strong> 
                                    <p>{{ comment.content }}</p>
                                </div>
                                <div class = "ml-auto"  v-if="canDeleteComment(comment,)" >
                                    <button class="btn btn-danger" @click="deleteComment(comment.id, comment.photoId, username, comment.userId)">Delete</button>
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
</style>