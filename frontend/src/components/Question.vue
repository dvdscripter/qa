<template>
<div v-if="question">
    <b-row align-v="center">
        <b-col cols="4">
            <div class="float-right">
                <b-row>
                    <b-button variant="click" @click="direction='PUT'">👍</b-button>
                </b-row>
                <b-row align-h="center" >
                    <!-- {{question.votes}} -->
                    <QuestionVote class="text-center" :votes="question.votes" :questionID="question.id" :direction="direction" @voted="question.votes += $event" />
                </b-row>
                <b-row>
                    <b-button variant="click" @click="direction ='DELETE'">👎</b-button>
                </b-row>
            </div>

        </b-col>
        <b-col  cols="4">
            <b-alert :show="hasError" variant="danger">{{error}}</b-alert>
            <!-- <h1 class="float-left">{{question.title}}</h1> -->
            <span class="question title">
              {{question.title}}  

            </span>
        </b-col>
        <b-col></b-col>
    </b-row>
    <b-row>
        <b-col cols="4"></b-col>
        <b-col cols="4" class="question content">
            <!-- <b-alert show variant="dark"> -->
                {{question.content}}
            <!-- </b-alert> -->
        </b-col>
    </b-row>
    <b-row >
        <b-col cols="4"></b-col>
        <b-col cols="4" >
          <div id="lastedit">
            <User class="" :userid="question.author" />
            <span  class="small">
              {{question.last_edit}}
            </span>

          </div>
            <b-button class="float-right" variant="success" @click="$router.push({name: 'CreateComment', params: {id: id}})">New Comment</b-button>

            <b-button @click="showCollapse = !showCollapse"
            :class="showCollapse ? 'collapsed' : null"
            aria-controls="collapse"
           :aria-expanded="showCollapse ? 'true' : 'false'"
            variant="primary" class="float-right">
            Edit
            </b-button>
            <b-collapse v-model="showCollapse" id="collapse">
                <b-form @submit.prevent="sendEdit">
                    <b-form-textarea v-model="editQuestion">
                    </b-form-textarea>
                    <b-button variant="primary" class="float-right" type="submit">Send</b-button>
                </b-form>
            </b-collapse>
        </b-col>

    </b-row>
    <b-row>
      <b-col cols="4"></b-col>
      <b-col cols="4" >
        <Comments :questionid="Number(id)" />
      </b-col>
    </b-row>
</div>
</template>

<script>
import QuestionVote from "./QuestionVote";
import Comments from "./Comments";
import User from "./User";
export default {
  name: "Question",
  components: {
    QuestionVote,
    Comments,
    User
  },
  data() {
    return {
      question: null,
      editQuestion: "",
      showCollapse: false,
      showCommentCollapse: false,
      direction: "",
      hasError: false,
      error: ""
    };
  },
  computed: {
    id() {
      return this.$route.params.id;
    }
  },
  methods: {
    getQuestion() {
      fetch(this.$APIENDPOINT + `/question/${this.id}`, {
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
          if (r["error"]) this.error = r["error"];
          else {
            this.question = r["result"];
            this.editQuestion = this.question.content;
          }
        })
        .catch(() => {
          throw Error("Cannot contact backend");
        });
    },
    sendEdit() {
      let newEdit = this.question;
      newEdit.content = this.editQuestion;
      fetch(this.$APIENDPOINT + `/question/${this.id}`, {
        method: "PUT",
        mode: "cors",
        body: JSON.stringify(newEdit),
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
            this.question = r["result"];
            this.editQuestion = this.question.content;
          }
        })
        .catch(e => {
          this.hasError = true;
          this.error = `Cannot contact backend: ${e.message}`;
        });
      this.showCollapse = false;
    }
  },
  mounted: function() {
    this.getQuestion();
  }
};
</script>

<style scoped>
textarea {
  margin-top: 3rem !important;
  margin-bottom: 1rem !important;
}
.question {
  word-wrap: break-word;
}
.title {
  font-size: 1.5em;
}
.content {
  border-top-style: solid;
  border-top-width: thin;
  /* margin-bottom: 1em; */
  padding-top: 1em;
}
#lastedit {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
</style>
