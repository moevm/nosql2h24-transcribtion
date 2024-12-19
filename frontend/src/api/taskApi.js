export const createTask = async (taskData) => {
    const response = await fetch('http://your-backend-url/api/tasks', {
      method: 'POST',
      body: taskData,
    });
    return response.json();
  };
  
  export const getTasks = async () => {
    const response = await fetch('http://your-backend-url/api/tasks');
    return response.json();
  };
  
  export const getTaskById = async (id) => {
    const response = await fetch(`http://your-backend-url/api/tasks/${id}`);
    return response.json();
  };
  
  export const updateTask = async (id, taskData) => {
    const response = await fetch(`http://your-backend-url/api/tasks/${id}`, {
      method: 'PUT',
      body: taskData,
    });
    return response.json();
  };
  
  export const deleteTask = async (id) => {
    const response = await fetch(`http://your-backend-url/api/tasks/${id}`, {
      method: 'DELETE',
    });
    return response.json();
  };