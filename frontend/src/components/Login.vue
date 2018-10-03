<template>
  <div class="login">
    <b-row>
      <b-col></b-col>
      <b-col>
        <b-alert :show="hasError" variant="danger">
          {{error}}
          </b-alert>
        <b-form @submit.prevent="sendLogin()">
            <b-form-input class="" required type="email" name="email" v-model="email" placeholder="Enter email" />
            <b-form-input class="" required type="password" name="password" v-model="password" placeholder="Enter Password" />
            <b-link class="float-right align-middle" :to="{name: 'CreateUser'}">New User?</b-link>
            <b-button class="float-right" type="submit" variant="primary">Submit</b-button>
        </b-form>
      </b-col>
      <b-col></b-col>
    </b-row>
  </div>
</template>

<script>
import isLogged from "../auth.js";
export default {
  name: "Login",
  data: function() {
    return {
      email: "",
      password: "",
      error: "",
      hasError: false
    };
  },
  methods: {
    sendLogin() {
      fetch(this.$APIENDPOINT + "/login", {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        body: JSON.stringify({ email: this.email, password: this.password }),
        headers: {
          "Content-Type": "application/json"
        }
      })
        .then(resp => {
          return resp.json();
        })
        .then(r => {
          if (r["error"]) {
            this.hasError = true;
            this.error = r["error"];
          } else {
            localStorage.setItem("token", r.result);
            this.$emit("login", isLogged());
            this.$router.push({ name: "Questions" });
          }
        })
        .catch(e => {
          this.hasError = true;
          this.error = `Cannot contact backend: ${e.message}`;
        });
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
a {
  margin-top: 0.4rem !important;
  /* color: #42b983; */
}
</style>
