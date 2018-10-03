<template>
  <b-container fluid>
    <div id="app">
      <NavBar :user="user" />
      <img src="./assets/logo.png" >
      <router-view @login="setLogin" /> 
    </div>
  </b-container>
</template>

<script>
import NavBar from "./components/NavBar";
import isLogged from "./auth";
export default {
  name: "app",
  components: {
    NavBar: NavBar
  },
  data() {
    return {
      user: ""
    };
  },
  methods: {
    setLogin(user) {
      this.user = user;
    }
  },
  mounted() {
    this.user = isLogged();
  }
};
</script>


<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  /* margin-top: 60px; */
}
</style>
