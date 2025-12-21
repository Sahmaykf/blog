import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import PostDetail from '../views/PostDetail.vue'
import Admin from '../views/Admin.vue'
import Login from '../views/Login.vue'
import UserProfile from '../views/UserProfile.vue'
import UserNetwork from '../views/UserNetwork.vue'
import UserSettings from '../views/UserSettings.vue'
import Search from '../views/Search.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/search',
    name: 'Search',
    component: Search
  },
  {
    path: '/post/:id',
    name: 'PostDetail',
    component: PostDetail
  },
  {
    path: '/admin',
    name: 'Admin',
    component: Admin,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/user/:id',
    name: 'UserProfile',
    component: UserProfile,
    props: true
  },
  {
    path: '/user/:id/network/:type',
    name: 'UserNetwork',
    component: UserNetwork,
    props: true
  },
  {
    path: '/settings',
    name: 'UserSettings',
    component: UserSettings,
    meta: { requiresAuth: true }
  }
  
  
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAuthenticated = localStorage.getItem('token')
  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
  } else {
    next()
  }
})

export default router
