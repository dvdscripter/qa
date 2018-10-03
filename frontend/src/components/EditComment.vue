<template>
    <div class="createcomment">
        <b-row>
        <b-col></b-col>
        <b-col>
            <b-alert :show="hasError" variant="danger">
            {{error}}
            </b-alert>
            <b-form @submit.prevent="sendComment">
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
  name: "EditComment",
  data() {
    return {
      hasError: false,
      error: null,
      content: "",
      comment: {}
    };
  },
  computed: {
    id() {
      return this.$route.params.id;
    },
    questionid() {
      return this.$route.params.questionid;
    }
  },
  mounted() {
    this.getComment();
  },
  methods: {
    sendComment() {
      this.comment.content = this.content;
      fetch(
        this.$APIENDPOINT + `/question/${this.questionid}/comments/${this.id}`,
        {
          method: "PUT",
          mode: "cors",
          body: JSON.stringify(this.comment),
          headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        }
      )
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
              params: { id: this.questionid }
            });
          }
        })
        .catch(e => {
          this.hasError = true;
          this.error = `Cannot contact backend: ${e.message}`;
        });
    },
    getComment() {
      fetch(
        this.$APIENDPOINT + `/question/${this.questionid}/comments/${this.id}`,
        {
          method: "GET",
          mode: "cors",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Bearer " + localStorage.getItem("token")
          }
        }
      )
        .then(resp => {
          return resp.json();
        })
        .then(r => {
          if (r["error"]) {
            this.error = r["error"];
            this.hasError = true;
          } else {
            this.comment = r.result;
            this.content = r.result.content;
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
