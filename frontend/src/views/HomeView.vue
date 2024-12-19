<template>
  <div>
    <h1>Welcome to the App</h1>
    <router-link v-if="!isAuthenticated" to="/register">Register</router-link>
    <router-link v-if="!isAuthenticated" to="/login">Login</router-link>
    <router-link v-if="isAuthenticated" to="/user-panel">User Panel</router-link>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useStore } from 'vuex';

export default {
  setup() {
    const store = useStore();
    const isAuthenticated = computed(() => store.getters.isAuthenticated);

    return {
      isAuthenticated,
    };
  },
};
</script>
