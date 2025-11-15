import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { env } from "@/config/env";
import { getSession } from "@/lib/supabase/client";
import { API_ERROR_MESSAGES } from "@/lib/constants";
import type { APIError } from "@/types/api";

/**
 * Create Axios instance with base configuration
 */
export const apiClient = axios.create({
  baseURL: env.api.baseUrl,
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

/**
 * Request interceptor to add auth token
 */
apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    try {
      const session = await getSession();
      if (session?.access_token) {
        config.headers.Authorization = `Bearer ${session.access_token}`;
      }
    } catch (error) {
      console.error("Error getting session:", error);
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * Response interceptor for error handling
 */
apiClient.interceptors.response.use(
  (response) => {
    // Return just the data for successful responses
    return response.data;
  },
  async (error: AxiosError<APIError>) => {
    // Handle different error cases
    if (error.response) {
      const status = error.response.status;
      const errorData = error.response.data;

      switch (status) {
        case 401:
          // Unauthorized - redirect to login
          if (typeof window !== "undefined") {
            window.location.href = "/login";
          }
          return Promise.reject({
            message: API_ERROR_MESSAGES.UNAUTHORIZED,
            status,
            ...errorData,
          });

        case 403:
          return Promise.reject({
            message: API_ERROR_MESSAGES.FORBIDDEN,
            status,
            ...errorData,
          });

        case 404:
          return Promise.reject({
            message: API_ERROR_MESSAGES.NOT_FOUND,
            status,
            ...errorData,
          });

        case 422:
          return Promise.reject({
            message: API_ERROR_MESSAGES.VALIDATION_ERROR,
            status,
            ...errorData,
          });

        case 500:
        case 502:
        case 503:
          return Promise.reject({
            message: API_ERROR_MESSAGES.SERVER_ERROR,
            status,
            ...errorData,
          });

        default:
          return Promise.reject({
            message: errorData?.message || "An error occurred",
            status,
            ...errorData,
          });
      }
    } else if (error.request) {
      // Network error
      return Promise.reject({
        message: API_ERROR_MESSAGES.NETWORK_ERROR,
        status: 0,
      });
    }

    return Promise.reject(error);
  }
);

/**
 * Helper function to handle API errors
 */
export function handleAPIError(error: any): string {
  if (error.message) {
    return error.message;
  }

  if (error.response?.data?.error) {
    return error.response.data.error;
  }

  return API_ERROR_MESSAGES.SERVER_ERROR;
}