// frontend/src/api/serverApi.js
const BASE_URL = 'http://localhost:8080';

export const getServers = async () => {
  const response = await fetch(`${BASE_URL}/servers`);
  return response.json();
};

export const getServerById = async (id) => {
  const response = await fetch(`${BASE_URL}/servers/${id}`);
  return response.json();
};

export const createServer = async (serverData) => {
  const response = await fetch(`${BASE_URL}/servers`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(serverData),
  });
  return response.json();
};

export const updateServer = async (id, serverData) => {
  const response = await fetch(`${BASE_URL}/servers/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(serverData),
  });
  return response.json();
};

export const patchServer = async (id, serverData) => {
  const response = await fetch(`${BASE_URL}/servers/${id}`, {
    method: 'PATCH',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(serverData),
  });
  return response.json();
};

export const deleteServer = async (id) => {
  const response = await fetch(`${BASE_URL}/servers/${id}`, {
    method: 'DELETE',
  });
  return response.json();
};

export const getServerCurrentJobs = async (id) => {
  const response = await fetch(`${BASE_URL}/servers/${id}/currentJobs`);
  return response.json();
};

export const getServerCompletedJobs = async (id) => {
  const response = await fetch(`${BASE_URL}/servers/${id}/completedJobs`);
  return response.json();
};

export const addJobToServer = async (serverId, jobId, jobData) => {
  const response = await fetch(`${BASE_URL}/servers/${serverId}/jobs/${jobId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(jobData),
  });
  return response.json();
};