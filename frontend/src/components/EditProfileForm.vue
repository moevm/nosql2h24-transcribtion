<template>
  <form @submit.prevent="handleEditProfile">
    <h2>Edit Profile</h2>

    <div>
      <label>Username:</label>
      <input type="text" v-model="username" />
    </div>

    <div>
      <label>Email:</label>
      <input type="email" v-model="email" />
    </div>

    <div>
      <label>Password:</label>
      <input type="password" v-model="password" />
    </div>

    <div>
      <label>Confirm Password:</label>
      <input type="password" v-model="confirmPassword" />
    </div>

    <button type="submit">Save Changes</button>
    <button type="button" @click="$emit('close')">Cancel</button>
    <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
  </form>
</template>

<script>
import { ref, defineProps } from 'vue';
import {useUserStore} from '../store/user';
import { updateUser } from '../api/userApi';

export default {
  props: {
    userData: Object,
  },
  setup() {
    const userStore = useUserStore();

    const username = ref('')
    const email = ref('');
    const password = ref('');
    const confirmPassword = ref('');
    const errorMessage = ref('');

    const handleEditProfile = () => {
      if (password.value !== confirmPassword.value) {
        errorMessage.value = 'Passwords do not match!';
        return;
      }

      errorMessage.value = '';

      const payload = {
        username: username.value ? username.value : userStore.username,
        email: email.value ? email.value : userStore.email,
        password_hash: password.value ? password.value : userStore.password_hash,
        permissions: 'user'
      };
      

      // Emit the updated data to the parent component
      try {
        updateUser(userStore.id, payload)

        userStore.username = payload.username;
        userStore.email = payload.email;
        userStore.password_hash = payload.password_hash;
      } catch (error) {
        errorMessage.value = 'Failed to update profile';
      }


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
  background-color: rgb(87, 57, 101);
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

