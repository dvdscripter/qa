<template>
  <div></div>
</template>


<script>
export default {
  name: "Logout",
  mounted: function() {
    localStorage.removeItem("token");
    this.$emit("login", "");
    this.$router.push({ name: "Login" });
  }
};
</script>
