<template>
    <b-table v-if="comments" stacked responsive stripped hover :fields="fields" :items="comments">
        <template slot="author" slot-scope="data">
            <User @user="setUsers($event)" :userid="data.item.author" />
        </template>
        <template slot="votes" slot-scope="data">
            <b-button variant="click" @click="Vote(1, data.index)">👍</b-button>
            {{data.item.votes}}
            <b-button variant="click" @click="Vote(-1, data.index)">👎</b-button>
        </template>
        <template slot="actions" slot-scope="data">
            <b-button @click="$router.push({name: 'EditComment', params: {questionid: data.item.question, id: data.item.id} })">Edit</b-button>
        </template>
    </b-table>
</template>

<script>
import User from "./User";
export default {
  name: "Comments",
  components: {
    User
  },
  props: {
    questionid: Number
  },
  data() {
    return {
      comments: [],
      users: [],
      fields: ["votes", "content", "author", "last_edit", "when", "actions"]
    };
  },
  mounted() {
    this.getComments();
  },
  methods: {
    setUsers(e) {
      this.users[e.id] = e;
    },
    Vote(direction, comment) {
      fetch(
        this.$APIENDPOINT +
          `/question/${this.questionid}/comments/${
            this.comments[comment].id
          }/vote`,
        {
          method: direction == 1 ? "PUT" : "DELETE",
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
          if (r["error"]) throw Error(r["error"]);
          else {
            this.comments[comment].votes += direction;
          }
        })
        .catch(e => {
          throw e;
        });
    },
    getComments() {
      fetch(this.$APIENDPOINT + `/question/${this.questionid}/comments`, {
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
            this.comments = r["result"];
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
</style>
