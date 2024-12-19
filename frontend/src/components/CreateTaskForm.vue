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
import {useUserStore} from '../store/user';


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

    const userStore = useUserStore();

    const handleFileUpload = (event) => {
      job.value.input_file = event.target.files[0];
    };

    const handleCreateJob = async () => {
      const userId = userStore.id; // Replace with actual user ID

      const formData = {
        title: job.value.title,
        status: job.value.status,
        source_language: job.value.source_language,
        file_format: job.value.file_format,
        description: job.value.description,
        input_file: job.value.input_file.name,
        output_file: job.value.output_file,
      }

      try {
        console.log(formData);

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
        alert('Failed to create job: ', error.message);
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
  background-color: #6b827a;
}
</style>