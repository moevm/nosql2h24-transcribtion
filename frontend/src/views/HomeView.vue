<template>
  <div>
    <h1>Welcome to the App</h1>
    <router-link v-if="!this.userStore.id" to="/register">Register</router-link>
    <router-link v-if="!this.userStore.id" to="/login">Login</router-link>
    <router-link v-if="this.userStore.id" to="/user-panel">User Panel</router-link>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';
import {useUserStore} from '../store/user';

export default {
  setup() {
    const store = useStore();
    const userStore = useUserStore();
    const isAuthenticated = computed(() => store.getters.isAuthenticated);

    return {
      isAuthenticated,
      userStore
    };
  },
};
</script>
