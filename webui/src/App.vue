<script setup>
import { RouterLink, RouterView } from 'vue-router'

</script>
<script>
export default {
	data() {
		return {
			errormsg: null,
			successmsg: null,
			detailedmsg: null,
			username: localStorage.getItem('username'),
			token: localStorage.getItem('token'),
			searchByUsername: "",
		};

	},
	methods: {
		async SearchUser() {
			if (this.searchByUsername === this.username) {
				this.errormsg = "You can't search yourself."
			} else if (this.searchByUsername === "") {
				this.errormsg = "Emtpy username field."
			} else {
				try {
					let response = await this.$axios.get("user/" + this.searchByUsername + "/profile", {
						headers: {
							Authorization: "Bearer " + localStorage.getItem("token")
						}
					})
					this.profile = response.data
					this.$router.push({ path: '/user/' + this.searchByUsername + '/view' })
				} catch (e) {
					if (e.response && e.response.status === 400) {
						this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
						this.detailedmsg = null;
					} else if (e.response && e.response.status === 500) {
						this.errormsg = "User does not exist on WASAPhoto.";
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
		<header class="navbar sticky-top shadow " style="background-color: #264653; color: #ffffff;">
			<div class="container-fluid">
				<!-- Logo all'estremo sinistro -->
				<Router-link to="/session">
					<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-9" style="background-color: #264653; color: #ffffff;">
						WASAPhoto
					</a>
   				</Router-link>
				<form class="d-flex mx-auto">
					<input class="form-control me-2" type="search" placeholder="Search" aria-label="Search">
					<button class="btn btn-light" type="submit">Search</button>
				</form>
        	</div>
		</header>

		<div class="container-fluid">
			<div class="row">
				<main class="col-md-9 ms-sm-auto col-lg-12 px-md-4">
					<RouterView />
				</main>
			</div>
		</div>

	</div>

</template>
