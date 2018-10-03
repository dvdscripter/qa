<template>
    <div class="createquestion">
    <b-row>
      <b-col></b-col>
      <b-col>
        <b-alert :show="hasError" variant="danger">
          {{error}}
        </b-alert>
        <b-form @submit.prevent="sendQuestion">
            <b-form-input required name="title" v-model="title" placeholder="Enter Title" />
            <b-form-textarea required v-model="content" />
            <b-button class="float-right" type="submit" variant="primary">Submit</b-button>
        </b-form>
      </b-col>
      <b-col></b-col>
    </b-row>
  </div>
</template>

<script>
export default {
  name: "CreateQuestion",
  data() {
    return {
      title: "",
      content: "",
      hasError: false,
      error: ""
    };
  },
  methods: {
    sendQuestion() {
      fetch(this.$APIENDPOINT + `/question`, {
        method: "POST",
        mode: "cors",
        body: JSON.stringify({ title: this.title, content: this.content }),
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
              params: { id: r.result.id }
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
