<template>
    <b-link v-if="user" :to="{name: 'Profile', params: {email: user.email}}">
      <b-img class="uimg" v-if="hasAvatar" :src="user.avatar" width="40" height="40"></b-img>
      <b-img class="uimg" v-if="!hasAvatar" blank width="40" height="40" blank-color="#777"></b-img>
    </b-link>
</template>

<script>
export default {
  name: "User",
  props: {
    userid: Number,
    email: String
  },
  data() {
    return {
      user: null,
      hasAvatar: false
    };
  },
  mounted() {
    if (this.userid !== undefined) this.getUser();
    else this.getUserByEmail();
  },
  methods: {
    getUser() {
      fetch(this.$APIENDPOINT + `/user/${this.userid}`, {
        method: "GET",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + localStorage.getItem("token")
        }
      })
        .then(resp => {
          return resp.json();
        })
        .then(r => {
          if (r["error"]) throw Error(r["error"]);
          else {
            this.user = r["result"];
            if (this.user.avatar) this.hasAvatar = true;
            this.$emit("user", this.user);
          }
        })
        .catch(e => {
          throw `Cannot contact backend: ${e.message}`;
        });
    },
    getUserByEmail() {
      fetch(this.$APIENDPOINT + `/user/${this.email}`, {
        method: "GET",
        mode: "cors",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Bearer " + localStorage.getItem("token")
        }
      })
        .then(resp => {
          return resp.json();
        })
        .then(r => {
          if (r["error"]) throw Error(r["error"]);
          else {
            this.user = r["result"];
            if (this.user.avatar) this.hasAvatar = true;
            this.$emit("user", this.user);
          }
        })
        .catch(e => {
          throw e;
        });
    }
  }
};
</script>

<style scoped>
.thumb-user {
  width: 40px;
  height: 40px;
}
.uimg {
  display: inline-block;
}
</style>
