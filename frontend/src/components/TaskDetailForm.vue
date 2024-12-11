<template>
  <div class="task-detail">
    <div class="header">
      <h2>{{ task.name }}</h2>
      <span class="status" :class="statusClass">{{ task.status }}</span>
    </div>

    <div class="details">
      <div>
        <label>From Language:</label>
        <p>{{ task.fromLanguage }}</p>
      </div>

      <div>
        <label>To Language:</label>
        <p>{{ task.toLanguage }}</p>
      </div>

      <div>
        <label>File Format:</label>
        <p>{{ task.fileFormat }}</p>
      </div>

      <div>
        <label>Description:</label>
        <p>{{ task.description }}</p>
      </div>
    </div>

    <button @click="downloadTask">Download Current Transcription Task</button>
  </div>
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    task: Object,
  },
  setup(props) {
    const downloadTask = () => {
      alert(`Downloading task: ${props.task.name}`);
      // Placeholder for download logic
    };

    const statusClass = computed(() => {
      return {
        'status-in-progress': props.task.status === 'In Progress',
        'status-completed': props.task.status === 'Completed',
        'status-pending': props.task.status === 'Pending',
      };
    });

    return { downloadTask, statusClass };
  },
};
</script>

<style scoped>
.task-detail {
  max-width: 600px;
  margin: auto;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background-color: #fff;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 2px solid #ccc;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

h2 {
  margin: 0;
}

.status {
  padding: 5px 10px;
  border-radius: 5px;
  color: #fff;
  font-weight: bold;
}

.status-in-progress {
  background-color: #ffa726;
}

.status-completed {
  background-color: #66bb6a;
}

.status-pending {
  background-color: #ef5350;
}

.details div {
  margin-bottom: 15px;
}

label {
  font-weight: bold;
  display: block;
}

p {
  margin: 5px 0;
}

button {
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
  width: 100%;
  font-size: 16px;
}

button:hover {
  background-color: #369870;
}
</style>
