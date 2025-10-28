import axios, { AxiosError } from "axios";
import type { InternalAxiosRequestConfig } from "axios";

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

// In-memory storage for access token
let accessToken: string | null = null;
export const setAccessToken = (token: string | null): void => {
  accessToken = token;
};

//Gets the current access token from memory.
export const getAccessToken = (): string | null => {
  return accessToken;
};

api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    if (accessToken && config.headers) {
      config.headers.Authorization = `Bearer ${accessToken}`;
    }
    return config;
  },
  (error: AxiosError) => Promise.reject(error)
);

//Shared promise to prevent multiple simultaneous refresh calls.
let refreshPromise: Promise<string | null> | null = null
async function refreshAccessToken(): Promise<string | null> {
  try {
    const { data } = await api.post("/users/refresh");
    const newToken = data.access_token as string;
    setAccessToken(newToken);
    return newToken;
  } catch (error) {
    // Refresh failed
    setAccessToken(null);
    return null;
  }
}
api.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _isRetryRequest?: boolean };
    
    // Only handle 401 errors
    if (!error.response || error.response.status !== 401) {
      return Promise.reject(error);
    }

    // Prevent infinite retry loop
    if (originalRequest._isRetryRequest) {
      return Promise.reject(error);
    }
    if (originalRequest.url?.endsWith('/users/refresh')) {
      return Promise.reject(error);
    }

    // Use shared promise to prevent concurrent refresh calls
    if (!refreshPromise) {
      refreshPromise = refreshAccessToken().finally(() => {
        refreshPromise = null;
      });
    }

    const newToken = await refreshPromise;

    if (newToken) {
      originalRequest._isRetryRequest = true;
      if (originalRequest.headers) {
        originalRequest.headers.Authorization = `Bearer ${newToken}`;
      }
      return api(originalRequest);
    }
    return Promise.reject(error);
  }
);

export default api;
