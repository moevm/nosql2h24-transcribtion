<template>
  <div class="user-panel">
    <h1>User Panel</h1>

    <div class="buttons">
      <button @click="openCreateTask">Create Task</button>
      <button @click="openEditProfile">Edit Profile</button>
      <button @click="openBillingForm">Billing</button>
      <button @click="openSearchPanel">Search for Tasks</button>
    </div>

    <table>
      <thead>
        <tr>
          <th>Task Name</th>
          <th>Status</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="task in tasks" :key="task.id" @click="selectTask(task)">
          <td>{{ task.name }}</td>
          <td>{{ task.status }}</td>
        </tr>
      </tbody>
    </table>

    <!-- Create Task Modal -->
    <div v-if="showCreateTask" class="modal-overlay">
      <div class="modal">
        <CreateTaskForm @create="handleCreateTask" @close="closeCreateTask" />
      </div>
    </div>

    <!-- Edit Profile Modal -->
    <div v-if="showEditProfile" class="modal-overlay">
      <div class="modal">
        <EditProfileForm :userData="userData" @save="handleSaveProfile" @close="closeEditProfile" />
      </div>
    </div>

    <!-- Task Detail Modal -->
    <div v-if="selectedTask" class="modal-overlay">
      <div class="modal">
        <TaskDetailForm :task="selectedTask" />
        <button class="close-button" @click="closeTaskDetail">Close</button>
      </div>
    </div>

    <!-- Billing Form Modal -->
    <div v-if="showBillingForm" class="modal-overlay">
      <div class="modal">
        <BillingForm @close="closeBillingForm" />
      </div>
    </div>

    <!-- Search Panel Modal -->
    <div v-if="showSearchPanel" class="modal-overlay">
      <div class="modal">
        <SearchPanel />
        <button class="close-button" @click="closeSearchPanel">Close</button>
      </div>
    </div>

  </div>
</template>

<script>
import { ref } from 'vue';
import EditProfileForm from '../components/EditProfileForm.vue';
import TaskDetailForm from '../components/TaskDetailForm.vue';
import CreateTaskForm from '../components/CreateTaskForm.vue';
import BillingForm from '../components/BillingForm.vue';
import SearchPanel from '../components/SearchPanel.vue';

export default {
  components: { EditProfileForm, TaskDetailForm, CreateTaskForm, BillingForm, SearchPanel },
  setup() {
    // Sample tasks data
    const tasks = ref([
      {
        id: 1,
        name: 'Design Homepage',
        status: 'In Progress',
        fromLanguage: 'English',
        toLanguage: 'Spanish',
        fileFormat: 'MP4',
        description: 'Transcribe homepage design meeting video.',
      },
      {
        id: 2,
        name: 'Fix Bug #102',
        status: 'Completed',
        fromLanguage: 'German',
        toLanguage: 'English',
        fileFormat: 'MP3',
        description: 'Transcribe bug report audio.',
      },
      {
        id: 3,
        name: 'Deploy to Production',
        status: 'Pending',
        fromLanguage: 'French',
        toLanguage: 'English',
        fileFormat: 'WAV',
        description: 'Transcribe deployment instructions.',
      },
    ]);

    const showEditProfile = ref(false);
    const userData = ref({
      username: 'JohnDoe',
      email: 'johndoe@example.com',
    });

    const selectedTask = ref(null);

    const selectTask = (task) => {
      selectedTask.value = task;
    };

    const closeTaskDetail = () => {
      selectedTask.value = null;
    };

    const showCreateTask = ref(false);

    const openCreateTask = () => {
      showCreateTask.value = true;
    };

    const closeCreateTask = () => {
      showCreateTask.value = false;
    };

    const showBillingForm = ref(false);
    const openBillingForm = () => (showBillingForm.value = true);
    const closeBillingForm = () => (showBillingForm.value = false);


    const showSearchPanel = ref(false);

    const openSearchPanel = () => {
      showSearchPanel.value = true;
    };

    const closeSearchPanel = () => {
      showSearchPanel.value = false;
    };


    // Button click handlers
    const createTask = () => {
      alert('Create Task button clicked');
    };

    const getAccess = () => {
      alert('Get Access button clicked');
    };

    const searchTasks = () => {
      alert('Search for Tasks button clicked');
    };

    const openEditProfile = () => {
      showEditProfile.value = true;
    };

    const closeEditProfile = () => {
      showEditProfile.value = false;
    };

    const handleSaveProfile = (updatedData) => {
      userData.value = { ...userData.value, ...updatedData };
      alert('Profile updated successfully!');
      closeEditProfile();
    };

    return {
      tasks,
      showEditProfile,
      userData,
      selectedTask,
      showCreateTask,
        openCreateTask,
        closeCreateTask,
    
        showBillingForm,
        openBillingForm,
        closeBillingForm,

        showSearchPanel,
      openSearchPanel,
      closeSearchPanel,

      selectTask,
      closeTaskDetail,
      createTask,
      getAccess,
      searchTasks,
      openEditProfile,
      closeEditProfile,
      handleSaveProfile,
    };
  },
};
</script>

<style scoped>
.user-panel {
  max-width: 800px;
  margin: auto;
  padding: 20px;
}

h1 {
  text-align: center;
}

.buttons {
  display: flex;
  justify-content: space-around;
  margin-bottom: 20px;
}

button {
  padding: 10px 20px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
  transition: background-color 0.3s;
}

button:hover {
  background-color: #369870;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th, td {
  border: 1px solid #ddd;
  padding: 10px;
  text-align: center;
}

th {
  background-color: #f4f4f4;
}

tr {
  cursor: pointer;
}

tr:hover {
  background-color: #f0f0f0;
}

/* Modal Styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.close-button {
  margin-top: 10px;
  background-color: #f44336;
}
</style>


