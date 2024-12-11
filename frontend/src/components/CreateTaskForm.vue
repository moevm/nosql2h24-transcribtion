<template>
  <form @submit.prevent="handleCreateTask" enctype="multipart/form-data">
    <h2>Create New Task</h2>

    <div>
      <label>Task Name:</label>
      <input type="text" v-model="taskName" required />
    </div>

    <div>
      <label>From Language:</label>
      <input type="text" v-model="fromLanguage" required />
    </div>

    <div>
      <label>To Language:</label>
      <input type="text" v-model="toLanguage" required />
    </div>

    <div>
      <label>File Format:</label>
      <input type="text" v-model="fileFormat" required />
    </div>

    <div>
      <label>Description:</label>
      <textarea v-model="description" required></textarea>
    </div>

    <div>
      <label>Upload File:</label>
      <input type="file" @change="handleFileUpload" required />
    </div>

    <button type="submit">Create Task</button>
    <button type="button" @click="$emit('close')">Cancel</button>
  </form>
</template>

<script>
import { ref } from 'vue';

export default {
  setup(_, { emit }) {
    const taskName = ref('');
    const fromLanguage = ref('');
    const toLanguage = ref('');
    const fileFormat = ref('');
    const description = ref('');
    const file = ref(null);

    const handleFileUpload = (event) => {
      file.value = event.target.files[0];
    };

    const handleCreateTask = () => {
      if (!file.value) {
        alert('Please upload a file');
        return;
      }

      const newTask = {
        id: Date.now(),
        name: taskName.value,
        fromLanguage: fromLanguage.value,
        toLanguage: toLanguage.value,
        fileFormat: fileFormat.value,
        description: description.value,
        file: file.value,
        status: 'Pending',
      };

      emit('create', newTask);
    };

    return {
      taskName,
      fromLanguage,
      toLanguage,
      fileFormat,
      description,
      handleFileUpload,
      handleCreateTask,
    };
  },
};
</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
  max-width: 400px;
  margin: auto;
  background-color: #f9f9f9;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

div {
  margin-bottom: 15px;
}

label {
  font-weight: bold;
  margin-bottom: 5px;
  display: block;
}

input[type="text"],
textarea,
input[type="file"] {
  padding: 8px;
  width: 100%;
  box-sizing: border-box;
}

button {
  padding: 10px;
  margin-top: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

button[type="button"] {
  background-color: #f44336;
}

button:hover {
  background-color: #369870;
}
</style>
