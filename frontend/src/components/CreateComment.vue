<template>
    <div class="createcomment">
        <b-row>
        <b-col></b-col>
        <b-col>
            <b-alert :show="hasError" variant="danger">
            {{error}}
            </b-alert>
            <b-form @submit.prevent="sendComment">
                <b-form-textarea required v-model="content" placeholder="Write a Comment" />
                <b-button class="float-right" type="submit" variant="primary">Submit</b-button>
            </b-form>
        </b-col>
        <b-col></b-col>
        </b-row>
    </div>
</template>

<script>
export default {
  name: "CreateComment",
  data() {
    return {
      hasError: false,
      error: null,
      content: ""
    };
  },
  computed: {
    id() {
      return this.$route.params.id;
    }
  },
  methods: {
    sendComment() {
      fetch(this.$APIENDPOINT + `/question/${this.id}/comments`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({ content: this.content }),
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
            this.$router.push({
              name: "Question",
              params: { id: this.id }
            });
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

<style>
</style>
