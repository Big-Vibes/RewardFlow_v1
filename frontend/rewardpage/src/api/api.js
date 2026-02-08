import axios from 'axios';

// Axios instance for backend API calls
// baseURL points to the Go API served on port 4000 under `/api` prefix
// Ensure backend `main.go` and router prefix match this baseURL.
const api = axios.create({
  baseURL: 'http://localhost:4000/api', // Backend API URL
  headers: {
    'Content-Type': 'application/json',
  },
});

// Attaches Authorization header automatically if token exists in localStorage
// The app stores the access token under key `token` (set by AuthContext.login)
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token){
        config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
});

// Global response interceptor: on 401 clear auth and redirect to login
api.interceptors.response.use(
    res => res,
    err => {
        if (err.response && err.response.status === 401){
            // Clear local auth state and force login
            localStorage.removeItem('token');
            localStorage.removeItem('role');
            window.location.href = '/login';
        }
        return Promise.reject(err);
    }
);

export default api;