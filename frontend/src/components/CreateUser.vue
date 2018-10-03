<template>
    <div class="CreateUser">
    <b-row>
      <b-col></b-col>
      <b-col>
        <b-alert variant="danger" :show="hasError">{{error}}</b-alert>
        <b-form @submit.prevent="sendCreateUser()">
            <b-form-input required v-model="nick" placeholder="Enter nickname" />
            <b-form-input required type="email" name="email" v-model="email" placeholder="Enter email" />
            <b-form-input required type="password" name="password" v-model="password" placeholder="Enter Password" />
            <b-form-file placeholder="Select a image for avatar" @change="fileSelect" accept="image/jpeg, image/png, image/gif"></b-form-file>
            <b-button class="float-right" type="submit" variant="primary">Submit</b-button>
        </b-form>
      </b-col>
      <b-col></b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "CreateUser",
  data() {
    return {
      error: null,
      email: "",
      password: "",
      nick: "",
      avatar: null,
      hasError: false
    };
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
    sendCreateUser() {
      fetch(this.$APIENDPOINT + "/user", {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        body: JSON.stringify({
          email: this.email,
          password: this.password,
          nick: this.nick,
          avatar: this.avatar
        }),
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
            this.$router.push({ name: "Login" });
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
