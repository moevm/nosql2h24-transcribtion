// frontend/src/api/userApi.js
const BASE_URL = 'http://your-backend-url/api';

export const getUsers = async () => {
  const response = await fetch(`${BASE_URL}/users`);
  return response.json();
};

export const getUserById = async (id) => {
  const response = await fetch(`${BASE_URL}/users/${id}`);
  return response.json();
};

export const createUser = async (userData) => {
  const response = await fetch(`${BASE_URL}/users`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });
  return response.json();
};

export const updateUser = async (id, userData) => {
  const response = await fetch(`${BASE_URL}/users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });
  return response.json();
};

export const patchUser = async (id, userData) => {
  const response = await fetch(`${BASE_URL}/users/${id}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(userData),
  });
  return response.json();
};

export const deleteUser = async (id) => {
  const response = await fetch(`${BASE_URL}/users/${id}`, {
    method: 'DELETE',
  });
  return response.json();
};

export const getUserJobs = async (id) => {
  const response = await fetch(`${BASE_URL}/users/${id}/jobs`);
  return response.json();
};

export const addUserJob = async (userId, formData) => {
  const response = await fetch(`${BASE_URL}/users/${userId}/jobs`, {
    method: 'POST',
    body: formData,
  });
  return response.json();
};

export const deleteUserJob = async (id, jobId) => {
  const response = await fetch(`${BASE_URL}/users/${id}/jobs/${jobId}`, {
    method: 'DELETE',
  });
  return response.json();
};

export const addPayment = async (id, paymentData) => {
  const response = await fetch(`${BASE_URL}/users/${id}/payments`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(paymentData),
  });
  return response.json();
};

export const deletePayment = async (id, paymentId) => {
  const response = await fetch(`${BASE_URL}/users/${id}/payments/${paymentId}`, {
    method: 'DELETE',
  });
  return response.json();
};