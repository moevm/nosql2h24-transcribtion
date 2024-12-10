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
import { ref } from 'vue';

export default {
  setup() {
    const email = ref('');
    const password = ref('');
    const errorMessage = ref('');

    const handleLogin = () => {
      if (!email.value || !password.value) {
        errorMessage.value = 'Please fill out all fields!';
        return;
      }
      errorMessage.value = '';
      console.log({ email: email.value, password: password.value });
      // Call API or handle login logic here
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
