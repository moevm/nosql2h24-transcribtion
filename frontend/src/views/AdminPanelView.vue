<template>
  <div class="admin-panel">
    <h1>Admin Panel</h1>
    <ServerForm />

    <button @click="loadServers">Refresh Servers</button>
    <h2>Servers</h2>
    <ul>
      <li v-for="server in servers" :key="server.id">
        <h3>{{ server.hostname }}</h3>
        <p>{{ server.description }}</p>
        <p>Status: {{ server.status }}</p>
        <p>Address: {{ server.address }}</p>
        <p>CPU: {{ server.cpu_info }}</p>
        <p>GPU: {{ server.gpu_info }}</p>
        <p>RAM: {{ server.ram_size_gb }} GB</p>
      </li>
    </ul>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { getServers } from '../api/serverApi';
import ServerForm from '../components/ServerForm.vue';

export default {
  components: { ServerForm },
  setup() {
    const servers = ref([]);

    const loadServers = async () => {

        servers.value = []
      try {
        const response = await getServers();
        console.log('SERVERS:', response);
        response.forEach((server) => {
            console.log(server)
          servers.value.push(server);
        });
      } catch (error) {
        alert('Failed to load servers');
      }
    };

    onMounted(loadServers);

    return {
      servers,
      loadServers,
    };
  },
};
</script>

<style scoped>
.admin-panel {
  max-width: 1000px;
  margin: auto;
  padding: 20px;
}
h1, h2 {
  text-align: center;
}
ul {
  display: flex;
  flex-wrap: wrap;
  list-style-type: none;
  padding: 0;
  gap: 10px;
}
li {
  border: 1px solid #ddd;
  padding: 10px;
  margin-bottom: 10px;
  flex: 1 1 calc(33.333% - 20px); /* Adjust the width as needed */
  box-sizing: border-box;
}
</style>
