<template>
  <form @submit.prevent="handleCreateTask" enctype="multipart/form-data">
    <h2>Create New Task</h2>
    <div>
      <label>Server Id:</label>
      <input type="text" v-model="serverId" required />
    </div>
    <div>
      <label>Job Id:</label>
      <input type="text" v-model="jobId" required />
    </div>
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
    const serverId = ref('');
    const file = ref(null);
    const jobId = ref('');

    const handleFileUpload = (event) => {
      file.value = event.target.files[0];
    };

    const handleCreateTask = async () => {
      if (!file.value) {
        alert('Please upload a file');
        return;
      }

      const formData = new FormData();
      formData.append('name', taskName.value);
      formData.append('fromLanguage', fromLanguage.value);
      formData.append('toLanguage', toLanguage.value);
      formData.append('fileFormat', fileFormat.value);
      formData.append('description', description.value);
      formData.append('file', file.value);

      console.log('http://localhost:5000/' + serverId.value + '/jobs/' + jobId.value)

      try {
        const response = await fetch('http://localhost:5000/' + serverId.value + '/jobs/' + jobId.value, {
          method: 'POST',
          body: formData,
        });

        if (!response.ok) {
          throw new Error('Failed to create task');
        }

        const newTask = await response.json();
        emit('create', newTask);
        alert('Task created successfully!');
      } catch (error) {
        alert(error.message);
      }
    };

    return {
      taskName,
      fromLanguage,
      toLanguage,
      fileFormat,
      description,
      serverId,
      jobId,
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
  background-color: #000000;
}
</style>
