import { ref } from 'vue'
import { defineStore } from 'pinia'

export const useUserStore = defineStore('userStore', () => {
  const id = ref(null);
  const username = ref('');
  const email = ref('');
  const password_hash = ref('');


  return { id, username, email, password_hash }
})
