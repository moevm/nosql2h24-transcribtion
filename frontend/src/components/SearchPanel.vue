<template>
  <div class="search-panel">
    <h2>Search Panel</h2>

    <!-- Search Box -->
    <div class="search-box">
      <div>
        <label>Source Language:</label>
        <input type="text" v-model="filters.source_language" placeholder="Source language" />
      </div>

      <div>
        <label>Title:</label>
        <input type="text" v-model="filters.title" placeholder="Title name" />
      </div>

      <div>
        <label>File format:</label>
        <input type="text" v-model="filters.file_format" placeholder="" />
      </div>

      <button @click="applyFilters">Search</button>
    </div>

    <!-- Sort By Button -->
    <div class="sort-section">
      <label>Sort By:</label>
      <select v-model="sortKey">
        <option value="title">Title</option>
        <option value="source_language">Source language</option>
      </select>
      <button @click="sortTasks">Sort</button>
    </div>

    <!-- Tasks Table -->
    <table>
      <thead>
        <tr>
          <th>Title</th>
          <th>Source language</th>
          <th>Created At</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="task in filteredTasks" :key="task.id">
          <td>{{ task.title }}</td>
          <td>{{ task.source_language }}</td>
          <td>{{ task.created_at }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useUserStore} from '../store/user';
import { getUserJobs } from '../api/userApi';

export default {
  setup() {
    // Sample tasks data
    const tasks = ref([]);
    const userStore = useUserStore();

    const filters = ref({
      source_language: '',
      title: '',
      file_format: '',
    });

    const loadTasks = async () => {
      tasks.value = await getUserJobs(userStore.id);
    };

    loadTasks();

    const sortKey = ref('name');
    const sortOrder = ref('asc');

    // Filtered and sorted tasks
    const filteredTasks = computed(() => {
      return tasks.value
        .filter((task) => {
          return (
            (!filters.value.source_language || task.source_language.toLowerCase().includes(filters.value.source_language.toLowerCase())) &&
            (!filters.value.title || task.title.toLowerCase().includes(filters.value.title.toLowerCase())) &&
            (!filters.value.file_format || task.title.toLowerCase().includes(filters.value.file_format.toLowerCase()))
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
      userStore,
      loadTasks
    };
  },
};
</script>

<style scoped>
.search-panel {
  max-width: 800px;
  margin: auto;
  padding: 20px;
  background-color: #452c44;
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
  background-color: #6b827a;
}

th, td {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: left;
  background-color: #6b827a;
}

th {
  background-color: #68af96;
}

tr:hover {
  background-color: #f0f0f0;
}
</style>
