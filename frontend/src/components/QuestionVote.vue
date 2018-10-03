<template>
    <div>{{votes}}</div>
</template>

<script>
export default {
  name: "QuestionVote",
  data() {
    return {};
  },
  props: {
    direction: String,
    questionID: Number,
    votes: Number
  },
  watch: {
    direction() {
      this.vote();
    }
  },
  methods: {
    vote() {
      fetch(this.$APIENDPOINT + `/question/${this.questionID}/vote`, {
        method: this.direction,
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
            this.$emit("voted", this.direction === "PUT" ? 1 : -1);
          }
        })
        .catch(() => {
          throw Error("Cannot contact backend");
        });
    }
  }
};
</script>
