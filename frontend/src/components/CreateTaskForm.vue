<template>
  <form @submit.prevent="handleCreateJob" enctype="multipart/form-data">
    <h2>Create New Job</h2>
    <div>
      <label>Title:</label>
      <input type="text" v-model="job.title" required />
    </div>
    <div>
      <label>Status:</label>
      <input type="text" v-model="job.status" required />
    </div>
    <div>
      <label>Source Language:</label>
      <input type="text" v-model="job.source_language" required />
    </div>
    <div>
      <label>File Format:</label>
      <input type="text" v-model="job.file_format" required />
    </div>
    <div>
      <label>Description:</label>
      <textarea v-model="job.description" required></textarea>
    </div>
    <div>
      <label>Upload Input File:</label>
      <input type="file" @change="handleFileUpload" required />
    </div>
    <div>
      <label>Output File:</label>
      <input type="text" v-model="job.output_file" required />
    </div>
    <button type="submit">Create Job</button>
  </form>
</template>

<script>
import { ref } from 'vue';
import { addUserJob } from '../api/userApi';

export default {
  setup() {
    const job = ref({
      title: '',
      status: 'pending',
      source_language: '',
      file_format: '',
      description: '',
      input_file: null,
      output_file: '',
    });

    const handleFileUpload = (event) => {
      job.value.input_file = event.target.files[0];
    };

    const handleCreateJob = async () => {
      const userId = 'user-id'; // Replace with actual user ID
      const formData = new FormData();
      formData.append('title', job.value.title);
      formData.append('status', job.value.status);
      formData.append('source_language', job.value.source_language);
      formData.append('file_format', job.value.file_format);
      formData.append('description', job.value.description);
      formData.append('input_file', job.value.input_file);
      formData.append('output_file', job.value.output_file);

      try {
        const newJob = await addUserJob(userId, formData);
        alert('Job created successfully!');
        // Optionally, reset the form or update the local state
        job.value = {
          title: '',
          status: 'pending',
          source_language: '',
          file_format: '',
          description: '',
          input_file: null,
          output_file: '',
        };
      } catch (error) {
        alert('Failed to create job');
      }
    };

    return {
      job,
      handleFileUpload,
      handleCreateJob,
    };
  },
};
</script>

<style scoped>
form {
  display: flex;
  flex-direction: column;
}
</style>