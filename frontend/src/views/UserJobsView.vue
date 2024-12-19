<template>
  <div>
    <h1>User Jobs</h1>
    <ul>
      <li v-for="job in jobs" :key="job.id">{{ job.title }}</li>
    </ul>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import { getUserJobs } from '../api/userApi';

export default {
  setup() {
    const jobs = ref([]);
    const userId = ref(null);

    onMounted(async () => {
      const route = useRoute();
      userId.value = route.params.id;
      jobs.value = await getUserJobs(userId.value);
    });

    return { jobs };
  },
};
</script>