<template>
  <form @submit.prevent="handleRegister">
    <h2>Register</h2>

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

    <div>
      <label>Confirm Password:</label>
      <input type="password" v-model="confirmPassword" required />
    </div>

    <button type="submit">Register</button>
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </form>
</template>

<script>
import { ref } from 'vue';

export default {
  setup() {
    const username = ref('');
    const email = ref('');
    const password = ref('');
    const confirmPassword = ref('');
    const errorMessage = ref('');

    const handleRegister = () => {
      if (password.value !== confirmPassword.value) {
        errorMessage.value = 'Passwords do not match!';
        return;
      }
      errorMessage.value = '';
      console.log({
        username: username.value,
        email: email.value,
        password: password.value,
      });
      // Call API or handle registration logic here
    };

    return { username, email, password, confirmPassword, handleRegister, errorMessage };
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
