<script>

export default {
    components: {},
    data: function () {
        return {
            errormsg: null,
            username: "",
            profile: {
                id: 0,
                username: "",
            },
        }
    },
    methods: {
        async doLogin() {
            if (this.username == "") {
                this.errormsg = "Username cannot be empty.";
            } else if (this.username.length > 20) {
              this.errormsg = "Username must be no more than 20 characters.";
            } else {
                try {
                    let response = await this.$axios.post("/session", { username: this.username })
                    this.profile = response.data
                    localStorage.setItem("token", this.profile.id);
                    localStorage.setItem("username", this.profile.username);
                    this.$router.push({ path: '/session' })
                } catch (e) {
                    if (e.response && e.response.status === 400) {
                        this.errormsg = "Form error, please check all fields and try again. If you think that this is an error, write an e-mail to us.";
                        this.detailedmsg = null;
                    } else if (e.response && e.response.status === 500) {
                        this.errormsg = "An internal error occurred. Please try again later.";
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

  <div class="container">
    <div class="row justify-content-center mt-5">
      <div class="col-md-6">
        <div class="card">
          <div class="card-body">
            <h2 class="card-title text-center mb-4">Login</h2>
            <form @submit.prevent="login">
              <div class="mb-3">
                <label for="username" class="form-label">Username:</label>
                <input v-model="username" type="text" class="form-control"  placeholder="Insert a username to log in." aria-label="Recipient's username"
             aria-describedby="basic-addon2">
              </div>
              <button class="btn btn-primary w-100" type="button" @click="doLogin">Login</button>
            </form>
          </div>
        </div>
      </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
    </div>
  </div>
</template>


<style>

</style>
