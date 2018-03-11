import Vue from 'vue'
import Router from 'vue-router'
import LandingPage from '@/components/LandingPage'
import QuestionPage from '@/components/QuestionPage'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'LandingPage',
      component: LandingPage
    },
    {
      path: '/test',
      name: 'QuestionPage',
      component: QuestionPage
    }
  ]
})
