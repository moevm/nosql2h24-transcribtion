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
import { createUser } from '../api/userApi';
import { useUserStore } from '../store/user';

export default {
  setup() {
    const username = ref('');
    const email = ref('');
    const password = ref('');
    const router = useRouter();
    const userStore = useUserStore();

    const handleRegister = async () => {
      try {
        const userData = {
          username: username.value,
          email: email.value,
          password_hash: password.value, // Assuming the backend hashes the password
        };
        const user = await createUser(userData);
        userStore.email = user.email;
        userStore.username = user.username;
        userStore.id = user.id;
        userStore.password_hash = user.password_hash;

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