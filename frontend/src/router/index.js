import { createRouter, createWebHistory } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import RegisterView from '../views/RegisterView.vue';
import LoginView from '../views/LoginView.vue';
import UserPanelView from '../views/UserPanelView.vue';
import UserJobsView from '../views/UserJobsView.vue';

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../views/AboutView.vue'),
    },
    { path: '/register', name: 'Register', component: RegisterView },
    { path: '/login', name: 'Login', component: LoginView },
    { path: '/user-panel', name: 'UserPanel', component: UserPanelView },
    { path: '/user/:id/jobs', name: 'UserJobs', component: UserJobsView },
  ],
})

export default router
