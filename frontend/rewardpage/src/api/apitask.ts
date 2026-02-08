// apitask.ts
// Helper functions for task-related API calls.
// Uses the shared `api` axios instance (`src/api/api.js`) which:
// - Automatically attaches Authorization header from localStorage token
// - Automatically redirects to login on 401 Unauthorized
// - Centralized error handling

import api from './api';

// CHANGE: Fetch today's daily tasks with cooldown and completion state
// Backend endpoint: GET /api/tasks/daily
// Expected response: {
//   tasks: [{ id, number, completed, completedAt }],
//   completedCount: number,
//   lastCompletedAt: ISO timestamp,
//   nextResetAt: ISO timestamp
// }
// CHANGE: Backend checks if tasks need reset (past midnight) and auto-resets them
const getDailyTasks = async (): Promise<any> => {
  try {
    const response = await api.get('/tasks/daily');
    return response.data;
  } catch (error: any) {
    console.error('Failed to fetch daily tasks:', error.response?.data || error.message);
    throw error;
  }
};

// CHANGE: Complete a task with backend cooldown validation
// Backend endpoint: POST /api/tasks/complete
// Body: { taskId, taskNumber }
// Backend validation:
// - Check if user completed 5 tasks today
// - Check if user is within 5-minute cooldown
// - Check if tasks need daily reset (past midnight)
// - Record completion timestamp in MongoDB
// Expected response: {
//   success: boolean,
//   task: { id, number, completed, completedAt },
//   completedCount: number,
//   nextResetAt: ISO timestamp,
//   cooldownUntil: ISO timestamp
// }
const completeTaskDaily = async (taskId: string): Promise<any> => {
  try {
    if (!taskId) {
      throw new Error('Task ID is required');
    }
    const response = await api.post('/tasks/complete', { taskId });
    return response.data;
  } catch (error: any) {
    console.error('Complete task failed:', error.response?.data || error.message);
    throw error;
  }
};

// CHANGE: Check cooldown status before allowing UI interaction
// Backend endpoint: GET /api/tasks/cooldown
// Returns: { isCooldownActive: boolean, remainingSeconds: number }
const checkCooldown = async (): Promise<any> => {
  try {
    const response = await api.get('/tasks/cooldown');
    return response.data;
  } catch (error: any) {
    console.error('Cooldown check failed:', error.response?.data || error.message);
    throw error;
  }
};

// Fetch all tasks for the logged-in user
// Backend endpoint: GET /api/tasks
// Expected response: array of task objects
const fetchTasks = async (): Promise<any[]> => {
  try {
    const response = await api.get('/tasks');
    return response.data || [];
  } catch (error: any) {
    console.error('Fetch tasks failed:', error.response?.data || error.message);
    throw error;
  }
};

// Create a new task
// Backend endpoint: POST /api/tasks
// Expected body: { title: string, completed: boolean }
// Expected response: { id, title, completed, createdAt }
const createTask = async (taskData: { title: string; completed?: boolean }): Promise<any> => {
  try {
    const response = await api.post('/tasks', taskData);
    return response.data;
  } catch (error: any) {
    console.error('Create task failed:', error.response?.data || error.message);
    throw error;
  }
};

// Complete a task by ID (legacy, use completeTaskDaily instead)
// Backend endpoint: POST /api/tasks/complete with taskId
const completeTask = async (taskId: string): Promise<any> => {
  try {
    if (!taskId) {
      throw new Error('Task ID is required');
    }
    const response = await api.post('/tasks/complete', { taskId, box: 'left' });
    return response.data;
  } catch (error: any) {
    console.error('Complete task failed:', error.response?.data || error.message);
    throw error;
  }
};

// Check in for daily streak
const checkIn = async (setStreak?: (streak: any) => void): Promise<any> => {
  try {
    const response = await api.post('/streak/update');
    const updatedUser = response.data;
    if (typeof setStreak === 'function' && updatedUser.streak) {
      setStreak(updatedUser.streak);
    }
    return updatedUser;
  } catch (error: any) {
    console.error('Check-in failed:', error.response?.data || error.message);
    throw error;
  }
};

// Load leaderboard rankings
const loadLeaderboard = async (setLeaderboard?: (data: any) => void): Promise<any> => {
  try {
    const response = await api.get('/leaderboard');
    const data = response.data;
    if (typeof setLeaderboard === 'function') {
      setLeaderboard(data);
    }
    return data;
  } catch (error: any) {
    console.error('Leaderboard fetch failed:', error.response?.data || error.message);
    throw error;
  }
};

export { getDailyTasks, checkCooldown, completeTaskDaily, completeTask, createTask, fetchTasks, checkIn, loadLeaderboard };
