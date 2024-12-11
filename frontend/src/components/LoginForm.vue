<template>
  <form @submit.prevent="handleLogin">
    <h2>Login</h2>

    <div>
      <label>Email:</label>
      <input type="email" v-model="email" required />
    </div>

    <div>
      <label>Password:</label>
      <input type="password" v-model="password" required />
    </div>

    <button type="submit">Login</button>
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </form>
</template>

<script>
import { useRouter } from 'vue-router';
import { ref } from 'vue';

export default {
  setup() {
    const email = ref('');
    const password = ref('');
    const errorMessage = ref('');
    const router = useRouter();

    const handleLogin = () => {
      if (!email.value || !password.value) {
        errorMessage.value = 'Please fill out all fields!';
        return;
      }
      errorMessage.value = '';

      // Simulate successful login (replace with real authentication logic)
      console.log({ email: email.value, password: password.value });
      router.push('/user-panel');
    };

    return { email, password, handleLogin, errorMessage };
  },
};

</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
  max-width: 300px;
  margin: auto;
}

div {
  margin-bottom: 10px;
}

label {
  margin-bottom: 5px;
  display: block;
}

input {
  padding: 8px;
}

button {
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>
