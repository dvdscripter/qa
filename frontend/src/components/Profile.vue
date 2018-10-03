<template>
    <div class="Profile">
      <template v-if="ownProfile">
        <b-row>
          <b-col></b-col>
          <b-col>
            <b-alert variant="danger" :show="hasError">{{error}}</b-alert>
            <b-form @submit.prevent="sendEditUser()">
                <b-form-input required v-model="nick" placeholder="Enter nickname"  />
                <b-form-input  type="password" name="password" v-model="password" placeholder="Enter Password" />
                <b-form-input  type="password" name="password" v-model="repPassword" placeholder="Enter Password Again" />
                <b-row>
                    <b-col>
                        <User class="float-left" :email="email" @user="setUser($event)" />
                    </b-col>
                    <b-col>
                        <b-form-file @change="fileSelect" placeholder="Select a image for avatar" v-model="avatar" accept="image/jpeg, image/png, image/gif"></b-form-file>
                    </b-col>
                </b-row>
                <b-button class="float-right" type="submit" variant="primary">Submit</b-button>
            </b-form>
          </b-col>
          <b-col></b-col>
        </b-row>
      </template>
      <div v-if="!ownProfile">
        <b-row >
          <b-col cols="4"></b-col>
          <b-col cols="4" class="justify-content-center" >
              <User :email="email" @user="setUser($event)" />
              <h1>
                {{nick}}
              </h1>
          </b-col>
          <b-col cols="4"></b-col>
        </b-row>
      </div>
  </div>
</template>

<script>
import User from "./User";
import isLogged from "../auth.js";
export default {
  name: "Profile",
  components: {
    User
  },
  data() {
    return {
      nick: "",
      password: null,
      repPassword: null,
      avatar: null,
      error: null,
      hasError: false,
      user: null
    };
  },
  computed: {
    email() {
      return this.$route.params.email;
    },
    ownProfile() {
      return this.email == isLogged();
    }
  },
  watch: {
    user() {
      this.nick = this.user.nick;
    }
  },
  methods: {
    fileSelect(e) {
      if (e) {
        let reader = new FileReader();
        reader.onload = event => {
          this.avatar = event.target.result;
        };
        reader.readAsDataURL(e.target.files[0]);
      }
    },
    setUser(u) {
      this.user = u;
    },
    sendEditUser() {
      let updateBody = { email: this.user.email, nick: this.nick };
      if (this.password) {
        if (this.password !== this.repPassword) {
          this.error = "Passwords aren't equal";
          this.hasError = true;
          return;
        }
        updateBody.password = this.password;
      }
      if (this.avatarb64()) updateBody.avatar = this.avatarb64();

      fetch(this.$APIENDPOINT + `/user/${this.user.id}`, {
        method: "PUT",
        mode: "cors",
        cache: "no-cache",
        body: JSON.stringify(updateBody),
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + localStorage.getItem("token")
        }
      })
        .then(resp => {
          return resp.json();
        })
        .then(r => {
          if (r["error"]) {
            this.error = r["error"];
            this.hasError = true;
          } else {
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

<style scoped>
</style>
