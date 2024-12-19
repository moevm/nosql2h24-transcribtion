<template>
  <div>
    <h1>Register</h1>
    <form @submit.prevent="handleRegister">
      <div>
        <label>Username:</label>
        <input type="text" v-model="username" required />
      </div>
      <div>
        <label>Email:</label>
        <input type="email" v-model="email" required />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" v-model="password" required />
      </div>
      <button type="submit">Register</button>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { createUser } from '../api/userApi';

export default {
  setup() {
    const username = ref('');
    const email = ref('');
    const password = ref('');
    const router = useRouter();
    const store = useStore();

    const handleRegister = async () => {
      try {
        const userData = {
          username: username.value,
          email: email.value,
          password_hash: password.value, // Assuming the backend hashes the password
        };
        const user = await createUser(userData);
        store.dispatch('login', user);
        alert('Registration successful!');
        router.push('/user-panel');
      } catch (error) {
        alert('Registration failed');
      }
    };

    return {
      username,
      email,
      password,
      handleRegister,
    };
  },
};
</script>