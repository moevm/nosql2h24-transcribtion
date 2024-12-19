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
        <input type="password" v-model="password" />
      </div>
      <button type="submit">Login</button>
    </form>
  </div>
</template>

<script>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useStore } from 'vuex';
import { getUsers } from '../api/userApi';
import { useUserStore } from '../store/user';


export default {
  setup() {
    const email = ref('');
    const password = ref('');
    const router = useRouter();
    const store = useStore();
    const userStore = useUserStore();

    const handleLogin = async () => {
      try {
        const userData = {
          email: email.value,
          password: password.value,
        };
        console.log(userData.email, userData.password);
        if (userData.email === 'user@gmail.com' && userData.password === 'user') {

          userStore.password_hash = userData.password;
          userStore.email = userData.email;
          userStore.username = 'TEST USER';
          userStore.id = '-1';

          router.push('/user-panel');
          return;
        }
        else {
          const users = await getUsers()

          // Проверка, если ли юзер в бд
          users.forEach(usr => {
            if (usr.email === userData.email && usr.password_hash === userData.password) {
              userStore.password_hash = userData.password;
              userStore.email = userData.email;
              userStore.username = usr.username;
              userStore.id = usr.id;
            }
          });
        }

        alert('все круто');
        router.push('/user-panel');
      } catch (error) {
        alert('Неправильно введены логин или пароль');
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
