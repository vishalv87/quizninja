import axios, { AxiosError, AxiosResponse, InternalAxiosRequestConfig } from "axios";
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

// Create base axios instance
const axiosInstance = axios.create({
  baseURL: getBaseURL(),
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Track rate limit state to short-circuit requests while rate-limited
let rateLimitedUntil = 0;

export function isRateLimited(): boolean {
  return Date.now() < rateLimitedUntil;
}

// Type the apiClient to return data directly
export const apiClient = axiosInstance as typeof axiosInstance & {
  get: <T = any>(url: string, config?: any) => Promise<T>;
  post: <T = any>(url: string, data?: any, config?: any) => Promise<T>;
  put: <T = any>(url: string, data?: any, config?: any) => Promise<T>;
  patch: <T = any>(url: string, data?: any, config?: any) => Promise<T>;
  delete: <T = any>(url: string, config?: any) => Promise<T>;
};

/**
 * Request interceptor to add auth token
 */
apiClient.interceptors.request.use(
  async (config: InternalAxiosRequestConfig) => {
    // Short-circuit requests while rate-limited to avoid wasting bandwidth
    if (isRateLimited()) {
      apiLogger.warn('Request blocked - client is rate-limited', { url: config.url });
      return Promise.reject({
        message: 'Too many requests. Please wait a moment.',
        status: 429,
      });
    }

    try {
      const session = await getSession();
      if (session?.access_token) {
        config.headers.Authorization = `Bearer ${session.access_token}`;
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
            ...errorData,
            message: errorData?.message || API_ERROR_MESSAGES.UNAUTHORIZED,
            status: errorData?.status || status,
          });

        case 403:
          return Promise.reject({
            ...errorData,
            message: errorData?.message || API_ERROR_MESSAGES.FORBIDDEN,
            status: errorData?.status || status,
          });

        case 404:
          return Promise.reject({
            ...errorData,
            message: errorData?.message || API_ERROR_MESSAGES.NOT_FOUND,
            status: errorData?.status || status,
          });

        case 422:
          return Promise.reject({
            ...errorData,
            message: errorData?.message || API_ERROR_MESSAGES.VALIDATION_ERROR,
            status: errorData?.status || status,
          });

        case 429: {
          // Rate limited — store the reset time to short-circuit future requests
          const retryAfter = (errorData as any)?.retry_after;
          const retryAfterHeader = error.response?.headers?.['retry-after'];
          if (retryAfter) {
            // Server returns unix timestamp in seconds
            rateLimitedUntil = retryAfter * 1000;
          } else if (retryAfterHeader) {
            rateLimitedUntil = parseInt(retryAfterHeader, 10) * 1000;
          } else {
            // Fallback: back off for 60 seconds
            rateLimitedUntil = Date.now() + 60000;
          }
          apiLogger.warn('Rate limited until', { resetAt: new Date(rateLimitedUntil).toISOString() });
          return Promise.reject({
            message: 'Too many requests. Please wait a moment.',
            status: 429,
            retryAfter: rateLimitedUntil,
          });
        }

        case 500:
        case 502:
        case 503:
          return Promise.reject({
            ...errorData,
            message: errorData?.message || API_ERROR_MESSAGES.SERVER_ERROR,
            status: errorData?.status || status,
          });

        default:
          return Promise.reject({
            ...errorData,
            message: errorData?.message || "An error occurred",
            status: errorData?.status || status,
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