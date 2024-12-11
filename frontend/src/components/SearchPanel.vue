<template>
  <div class="search-panel">
    <h2>Search Panel</h2>

    <!-- Search Box -->
    <div class="search-box">
      <div>
        <label>Task Type:</label>
        <input type="text" v-model="filters.taskType" placeholder="e.g., Transcription" />
      </div>

      <div>
        <label>Attribute:</label>
        <input type="text" v-model="filters.attribute" placeholder="e.g., Audio" />
      </div>

      <div>
        <label>Task Name:</label>
        <input type="text" v-model="filters.taskName" placeholder="e.g., Meeting Notes" />
      </div>

      <div>
        <label>Text:</label>
        <input type="text" v-model="filters.text" placeholder="Enter keyword" />
      </div>

      <button @click="applyFilters">Search</button>
    </div>

    <!-- Sort By Button -->
    <div class="sort-section">
      <label>Sort By:</label>
      <select v-model="sortKey">
        <option value="name">Task Name</option>
        <option value="type">Task Type</option>
        <option value="createdAt">Created At</option>
      </select>
      <button @click="sortTasks">Sort</button>
    </div>

    <!-- Tasks Table -->
    <table>
      <thead>
        <tr>
          <th>Task Name</th>
          <th>Task Type</th>
          <th>Created At</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="task in filteredTasks" :key="task.id">
          <td>{{ task.name }}</td>
          <td>{{ task.type }}</td>
          <td>{{ task.createdAt }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { ref, computed } from 'vue';

export default {
  setup() {
    // Sample tasks data
    const tasks = ref([
      { id: 1, name: 'Transcribe Meeting', type: 'Audio', createdAt: '2024-04-01' },
      { id: 2, name: 'Translate Report', type: 'Document', createdAt: '2024-05-10' },
      { id: 3, name: 'Subtitles for Video', type: 'Video', createdAt: '2024-06-15' },
    ]);

    const filters = ref({
      taskType: '',
      attribute: '',
      taskName: '',
      text: '',
    });

    const sortKey = ref('name');
    const sortOrder = ref('asc');

    // Filtered and sorted tasks
    const filteredTasks = computed(() => {
      return tasks.value
        .filter((task) => {
          return (
            (!filters.value.taskType || task.type.toLowerCase().includes(filters.value.taskType.toLowerCase())) &&
            (!filters.value.taskName || task.name.toLowerCase().includes(filters.value.taskName.toLowerCase())) &&
            (!filters.value.text || task.name.toLowerCase().includes(filters.value.text.toLowerCase()))
          );
        })
        .sort((a, b) => {
          let modifier = sortOrder.value === 'asc' ? 1 : -1;
          if (a[sortKey.value] < b[sortKey.value]) return -1 * modifier;
          if (a[sortKey.value] > b[sortKey.value]) return 1 * modifier;
          return 0;
        });
    });

    const applyFilters = () => {
      console.log('Filters applied:', filters.value);
    };

    const sortTasks = () => {
      sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc';
    };

    return {
      filters,
      sortKey,
      sortTasks,
      applyFilters,
      filteredTasks,
    };
  },
};
</script>

<style scoped>
.search-panel {
  max-width: 800px;
  margin: auto;
  padding: 20px;
}

h2 {
  text-align: center;
}

/* Search Box Styles */
.search-box {
  display: flex;
  flex-wrap: wrap;
  gap: 15px;
  margin-bottom: 20px;
}

.search-box div {
  flex: 1 1 calc(50% - 10px);
}

.search-box label {
  display: block;
  font-weight: bold;
  margin-bottom: 5px;
}

.search-box input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

.search-box button {
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
  flex: 1 1 100%;
}

.search-box button:hover {
  background-color: #369870;
}

/* Sort Section Styles */
.sort-section {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
  align-items: center;
}

.sort-section label {
  font-weight: bold;
}

.sort-section select {
  padding: 8px;
}

.sort-section button {
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
}

.sort-section button:hover {
  background-color: #369870;
}

/* Table Styles */
table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: left;
}

th {
  background-color: #f4f4f4;
}

tr:hover {
  background-color: #f0f0f0;
}
</style>
