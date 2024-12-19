<template>
  <div>
    <h1>Login</h1>
    <form @submit.prevent="handleLogin">
      <div>
        <label>Email:</label>
        <input type="email" v-model="email" required />
      </div>
      <div>
        <label>Password:</label>
        <input type="password" v-model="password" required />
      </div>
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { loginUser } from '../api/userApi';

export default {
  setup() {
    const email = ref('');
    const password = ref('');
    const router = useRouter();
    const store = useStore();

    const handleLogin = async () => {
      try {
        const userData = {
          email: email.value,
          password: password.value,
        };
        console.log(userData.email, userData.password);
        if (userData.email === 'user@gmail.com' && userData.password === 'user') {
          store.dispatch('login', { email: userData.email });
          router.push('/user-panel');
          return;
        }

        const user = await loginUser(userData);
        store.dispatch('login', user);
        alert('Login successful!');
        router.push('/user-panel');
      } catch (error) {
        alert('Login failed');
      }
    };

    return {
      email,
      password,
      handleLogin,
    };
  },
};
</script>
