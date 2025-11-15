import axios, { AxiosError, InternalAxiosRequestConfig } from "axios";
import { getSession } from "@/lib/supabase/client";
import { API_ERROR_MESSAGES } from "@/lib/constants";
import { apiLogger } from "@/lib/logger";
import type { APIError } from "@/types/api";

/**
 * Create Axios instance with base configuration
 * Using a factory function to defer environment variable access
 */
function getBaseURL(): string {
  // Fallback to localhost if env var not available
  return process.env.NEXT_PUBLIC_API_BASE_URL || "http://127.0.0.1:8080/api/v1";
}

export const apiClient = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

apiLogger.info('API client initialized', {
  baseURL: getBaseURL(),
  timeout: 30000,
});

/**
 * Request interceptor to add auth token
 */
apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    apiLogger.debug('API request starting', {
      method: config.method,
      url: config.url,
    });

    try {
      const session = await getSession();
      if (session?.access_token) {
        config.headers.Authorization = `Bearer ${session.access_token}`;
        apiLogger.debug('Added auth token to request');
      } else {
        apiLogger.debug('No session token available');
      }
    } catch (error) {
      apiLogger.error("Error getting session for API request", error);
    }
    return config;
  },
  (error) => {
    apiLogger.error('Request interceptor error', error);
    return Promise.reject(error);
  }
);

/**
 * Response interceptor for error handling
 */
apiClient.interceptors.response.use(
  (response) => {
    apiLogger.debug('API response received', {
      status: response.status,
      url: response.config.url,
    });
    // Return just the data for successful responses
    return response.data;
  },
  async (error: AxiosError<APIError>) => {
    // Handle different error cases
    if (error.response) {
      const status = error.response.status;
      const errorData = error.response.data;

      apiLogger.error('API response error', {
        status,
        url: error.config?.url,
        errorData,
      });

      switch (status) {
        case 401:
          apiLogger.warn('Unauthorized - redirecting to login');
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
      apiLogger.error('Network error - no response received', {
        url: error.config?.url,
      });
      return Promise.reject({
        message: API_ERROR_MESSAGES.NETWORK_ERROR,
        status: 0,
      });
    }

    apiLogger.error('Unknown API error', error);
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