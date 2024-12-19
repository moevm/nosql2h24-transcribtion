<template>
  <form @submit.prevent="handleSubmit">
    <h2>Create Server</h2>
    <div>
      <label>Hostname:</label>
      <input type="text" v-model="server.hostname" required />
    </div>
    <div>
      <label>Address:</label>
      <input type="text" v-model="server.address" required />
    </div>
    <div>
      <label>Description:</label>
      <textarea v-model="server.description" required></textarea>
    </div>
    <div>
      <label>Status:</label>
      <input type="text" v-model="server.status" required />
    </div>
    <div>
      <label>CPU Info:</label>
      <input type="text" v-model="server.cpu_info" required />
    </div>
    <div>
      <label>GPU Info:</label>
      <input type="text" v-model="server.gpu_info" required />
    </div>
    <div>
      <label>RAM Size (GB):</label>
      <input type="number" v-model="server.ram_size_gb" required />
    </div>
    <button type="submit">Create Server</button>
  </form>
</template>

<script>
import { ref } from 'vue';
import { createServer } from '../api/serverApi';

export default {
  setup() {
    const server = ref({
      hostname: '',
      address: '',
      description: '',
      status: 'active',
      cpu_info: '',
      gpu_info: '',
      ram_size_gb: 0,
    });

    const handleSubmit = async () => {
      try {
        await createServer(server.value);
        alert('Server created successfully!');
        // Optionally, reset the form
        server.value = {
          hostname: '',
          address: '',
          description: '',
          status: 'active',
          cpu_info: '',
          gpu_info: '',
          ram_size_gb: 0,
        };
      } catch (error) {
        alert('Failed to create server');
      }
    };

    return {
      server,
      handleSubmit,
    };
  },
};
</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
}
form div {
  margin-bottom: 10px;
}
form label {
  font-weight: bold;
}
form input, form textarea {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}
form button {
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
}
form button:hover {
  background-color: #369870;
}
</style>