<script setup>
import { RouterLink, RouterView } from 'vue-router'

</script>
<script>
export default {
	data() {
		return {
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
