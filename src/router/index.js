import Vue from 'vue'
import Router from 'vue-router'
import LandingPage from '@/components/LandingPage'
import QuestionPage from '@/components/QuestionPage'
import store from './../store'
Vue.use(Router)

function asdf(to, from, next)
{
  console.log(to.params.room);
  console.log(store.state.room);
  if(to.params.room != store.state.room.Code)
  {
    console.log(store.state.room.Code);
    console.log(to.params.room);
    console.log("fdsfad");
    next('/');
  }
  else{
    console.log(store.state.room);
    next();
  }
}

export default new Router({
  routes: [
    {
      path: '/',
      name: 'LandingPage',
      component: LandingPage
    },
    {
      path: '/join/:room',
      name: 'Join',
      component: QuestionPage,
      beforeEnter: asdf
    },
  ]
})
