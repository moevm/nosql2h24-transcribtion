<template>
  <form @submit.prevent="handleEditProfile">
    <h2>Edit Profile</h2>

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

    <button type="submit">Save Changes</button>
    <button type="button" @click="$emit('close')">Cancel</button>
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </form>
</template>

<script>
import { ref, defineProps } from 'vue';

export default {
  props: {
    userData: Object,
  },
  setup(props, { emit }) {
    const username = ref(props.userData.username);
    const email = ref(props.userData.email);
    const password = ref('');
    const confirmPassword = ref('');
    const errorMessage = ref('');

    const handleEditProfile = () => {
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

      // Emit the updated data to the parent component
      emit('save', {
        username: username.value,
        email: email.value,
        password: password.value,
      });
    };

    return { username, email, password, confirmPassword, handleEditProfile, errorMessage };
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
  margin-top: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
}

button[type="button"] {
  background-color: #f44336;
}

.error {
  color: red;
  margin-top: 10px;
}
</style>

