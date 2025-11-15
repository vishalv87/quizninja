/**
 * Environment variable validation and configuration
 */

function getEnvVar(key: string): string {
  const value = process.env[key];
  if (!value) {
    throw new Error(`Missing environment variable: ${key}`);
  }
  return value;
}

function getOptionalEnvVar(key: string, defaultValue: string = ""): string {
  return process.env[key] || defaultValue;
}

// Use getters to defer environment variable access until actually used
export const env = {
  // Supabase
  supabase: {
    get url() {
      return getEnvVar("NEXT_PUBLIC_SUPABASE_URL");
    },
    get anonKey() {
      return getEnvVar("NEXT_PUBLIC_SUPABASE_ANON_KEY");
    },
  },

  // API
  api: {
    get baseUrl() {
      return getEnvVar("NEXT_PUBLIC_API_BASE_URL");
    },
  },

  // App
  app: {
    get url() {
      return getEnvVar("NEXT_PUBLIC_APP_URL");
    },
    get name() {
      return getOptionalEnvVar("NEXT_PUBLIC_APP_NAME", "QuizNinja");
    },
  },

  // Optional
  sentry: {
    get dsn() {
      return getOptionalEnvVar("NEXT_PUBLIC_SENTRY_DSN");
    },
  },
  analytics: {
    get gaTrackingId() {
      return getOptionalEnvVar("NEXT_PUBLIC_GA_TRACKING_ID");
    },
  },
} as const;

// Validate environment on module load
export function validateEnv() {
  try {
    const required = [
      "NEXT_PUBLIC_SUPABASE_URL",
      "NEXT_PUBLIC_SUPABASE_ANON_KEY",
      "NEXT_PUBLIC_API_BASE_URL",
      "NEXT_PUBLIC_APP_URL",
    ];

    const missing = required.filter((key) => !process.env[key]);

    if (missing.length > 0) {
      throw new Error(
        `Missing required environment variables:\n${missing.map((k) => `  - ${k}`).join("\n")}`
      );
    }

    console.log("✅ Environment variables validated successfully");
  } catch (error) {
    console.error("❌ Environment validation failed:", error);
    throw error;
  }
}
