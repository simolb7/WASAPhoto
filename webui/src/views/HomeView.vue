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
	},
	mounted() {
	}
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Homepage</h1>
			<div class="btn-toolbar mb-2 mb-md-1">
				<div class="btn-group me-2">
					<input type="file" accept="image/*" class="btn btn-outline-primary" @change="uploadFile" ref="file">
					<button class="btn btn-success" @click="submitFile">Upload</button>
					<button class="btn btn-primary mx-5" type="button" @click="ViewProfile">Profilo</button>
				</div>
			</div>
		</div>

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>


