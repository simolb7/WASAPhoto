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

		};

	},
	methods: {
	},
	mounted() {
	}
}

</script>

<template>
	<div>
		<header class="navbar sticky-top shadow " style="background-color: #264653; color: #ffffff; height: 60px;">
			<div class="container-fluid">
				<!-- Logo all'estremo sinistro -->
				<Router-link to="/session">
					<a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-9" style="background-color: #264653; color: #ffffff;">
						WASAPhoto
					</a>
   				</Router-link>
        	</div>
		</header>

		<div class="col-md-9 ms-sm-auto col-lg-12 px-md-4 mt-4" > 
			<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
		</div>
	
		
		<div class="container-fluid">
			<div class="row">
				<main class="col-md-9 ms-sm-auto col-lg-12 px-md-4">
					<RouterView />
				</main>
			</div>
		</div>
		

	</div>

</template>
