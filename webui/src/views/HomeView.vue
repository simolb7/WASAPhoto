<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			username: localStorage.getItem('username'),
			token: localStorage.getItem('token'),
			images: null,
			profile: {
				requestId: 0,
				id: 0,
				username: "",
				followersCount: 0,
				followingCount: 0,
				photoCount: 0,
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
						usernameow: "",
                        
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
		async ViewProfile() {
          this.$router.push({ path: '/user/' + this.username + '/profile' })
        },
		async uploadFile() {
			this.images = this.$refs.file.files[0]
		},
		async submitFile() {
			if (this.images === null) {
				this.errormsg = "Please select a file."
			} else {
				try {
					let response = await this.$axios.post("/user/" + this.username + "/photo", this.images, {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}
					})
					this.profile = response.data
					this.images = null
					this.successmsg = "Photo uploaded successfully."
					this.$router.replace({ path: '/user/' + this.username + '/profile' });
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
		async streamFollowers() {
            try {
                let response = await this.$axios.get("/user/" + this.username + "/photos", {
                    headers: {
                        Authorization: "Bearer " + localStorage.getItem("token")
                    }
                })
                this.photoList = response.data
                
                for (let i = 0; i < this.photoList.photos.length; i++) {
                    this.photoList.photos[i].file = 'data:image/jpg;base64,' + this.photoList.photos[i].file;
                
                    let likestatus = await this.getLikeStatus(localStorage.getItem('username'), this.photoList.photos[i].id);

				
					let usow = await this.getusername(this.photoList.photos[i].userId);
					
					this.photoList.photos[i].usernameow = usow;
			

                   
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
                        this.errormsg = " Your following haven't post any photos yet";
                        this.detailedmsg = null;
                    }
                }
            }
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


		formatDateTime(dateTime) {
            //timeZoneName: 'short' 
            const options = { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit'};
            const formattedDateTime = new Date(dateTime).toLocaleString('en-US', options);
            return formattedDateTime;
        },
		async sendComment(photoid, comment, ownerid) {
            if (comment === "") {
                this.errormsg = "Emtpy comment field."
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
                let response = await this.$axios.get("/user/" + this.username  + "/photo/" + photoid + "/comment", {
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
		async refresh() {
            await this.streamFollowers();
            
        },
	},
	mounted() {
		this.streamFollowers()
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Homepage of {{ this.username }} </h1>
			<div class="btn-toolbar mb-2 mb-md-1">
				<div class="btn-group me-2">
					<input type="file" accept="image/*" class="btn btn-outline-primary" @change="uploadFile" ref="file">
					<button class="btn btn-success" @click="submitFile">Upload</button>
					<button class="btn btn-primary mx-5" type="button" @click="ViewProfile">Profile</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>

	<div class="row" v-if="photoList && photoList.photos && photoList.photos.length > 0">
        <div class="col-md-4" v-for="photo in photoList.photos" :key="photo.id">
            <div class="card mb-4 shadow-sm">
                <img class="card-img-top" :src=photo.file alt="Card image cap">
                <div class="card-body">
					<RouterLink  v-if="photo.usernameow" :to="'/user/' + photo.usernameow + '/view'" class="nav-link">
							<button type="button" class="btn btn-outline-primary expanded-button">@{{photo.usernameow}}</button>
						</RouterLink>
                   
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
                                    <strong>{{ comment.username }}</strong> 
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


