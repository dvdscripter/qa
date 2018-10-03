<template>
<div class='questions'>
  <b-row>
    <b-col>
      <b-link :to="{name: 'CreateQuestion'}">New Question</b-link>
    </b-col>
  </b-row>
  <b-row>
    <b-table v-if="questions" responsive stripped hover :fields="fields" :items="questions">
        <template slot="title" slot-scope="data">
            <b-link :to="{name: 'Question', params: {id: data.item.id}}">
                {{data.value}}
            </b-link>
        </template>
        <template slot="author" slot-scope="data">
          <User v-if="questions" :userid="data.item.author" />
        </template>
    </b-table>
  </b-row>
</div>
</template>

<script>
import User from "./User";
export default {
  name: "Questions",
  components: {
    User
  },
  data() {
    return {
      questions: [],
      fields: ["votes", "title", "author", "last_edit"]
    };
  },
  methods: {
    getQuestions() {
      fetch(this.$APIENDPOINT + "/question", {
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
          if (r["error"]) throw Error(["error"]);
          else this.questions = r["result"];
        })
        .catch(() => {
          throw Error("Cannot contact backend");
        });
    }
  },
  mounted: function() {
    this.getQuestions();
  }
};
</script>

<style scoped>
.table td {
  text-align: center;
}
</style>
