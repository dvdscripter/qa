import Vue from "vue";
import Router from "vue-router";
import Login from "./components/Login.vue";
import Logout from "./components/Logout.vue";
import Questions from "./components/Questions.vue";
import Question from "./components/Question.vue";
import CreateQuestion from "./components/CreateQuestion.vue";
import CreateUser from "./components/CreateUser.vue";
import CreateComment from "./components/CreateComment.vue";
import EditComment from "./components/EditComment.vue";
import Profile from "./components/Profile.vue";
import isLogged from "./auth";

Vue.use(Router);
Vue.prototype.$APIENDPOINT = "http://localhost:8000";

const router = new Router({
  routes: [
    {
      path: "/",
      name: "home",
      redirect: { name: "Questions" }
    },
    {
      path: "/login",
      name: "Login",
      component: Login,
      meta: { public: true }
    },
    {
      path: "/createuser",
      name: "CreateUser",
      component: CreateUser,
      meta: { public: true }
    },
    {
      path: "/logout",
      name: "Logout",
      component: Logout
    },
    {
      path: "/questions",
      name: "Questions",
      component: Questions
    },
    {
      path: "/question/:id",
      name: "Question",
      component: Question
    },
    {
      path: "/createquestion",
      name: "CreateQuestion",
      component: CreateQuestion
    },
    {
      path: "/question/:id/comments",
      name: "CreateComment",
      component: CreateComment
    },
    {
      path: "/question/:questionid/comments/:id",
      name: "EditComment",
      component: EditComment
    },
    {
      path: "/profile/:email",
      name: "Profile",
      component: Profile
    }
  ]
});

router.beforeEach((to, from, next) => {
  if (to.matched.some(record => record.meta.public)) {
    if (isLogged() && (to.name == "Login" || to.name == "CreateUser"))
      next({ name: "Questions" });
    else next();
  } else {
    if (!isLogged()) next({ name: "Login" });
    else next();
  }
});

export default router;
