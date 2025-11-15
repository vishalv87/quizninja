import { createClientComponentClient } from "@supabase/auth-helpers-nextjs";
import { supabaseLogger } from "@/lib/logger";

// Access environment variables directly to avoid module initialization timing issues
const SUPABASE_URL = process.env.NEXT_PUBLIC_SUPABASE_URL;
const SUPABASE_ANON_KEY = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY;

// Validate environment variables
if (!SUPABASE_URL) {
  throw new Error(
    "Missing NEXT_PUBLIC_SUPABASE_URL environment variable. Please check your .env.local file."
  );
}

if (!SUPABASE_ANON_KEY) {
  throw new Error(
    "Missing NEXT_PUBLIC_SUPABASE_ANON_KEY environment variable. Please check your .env.local file."
  );
}

supabaseLogger.info("Initializing Supabase client (cookie-based)", {
  url: SUPABASE_URL,
  hasKey: !!SUPABASE_ANON_KEY,
});

// Use cookie-based client for Next.js App Router
// This ensures middleware can read the session from cookies
export const supabase = createClientComponentClient();

supabaseLogger.info("Supabase client initialized successfully");

/**
 * Get the current session
 */
export async function getSession() {
  supabaseLogger.debug("Getting current session");
  const { data, error } = await supabase.auth.getSession();
  if (error) {
    supabaseLogger.error("Failed to get session", error);
    throw error;
  }
  supabaseLogger.debug("Session retrieved", { hasSession: !!data.session });
  return data.session;
}

/**
 * Get the current user
 */
export async function getUser() {
  supabaseLogger.debug("Getting current user");
  const { data, error } = await supabase.auth.getUser();
  if (error) {
    supabaseLogger.error("Failed to get user", error);
    throw error;
  }
  supabaseLogger.debug("User retrieved", { userId: data.user?.id });
  return data.user;
}

/**
 * Sign in with email and password
 */
export async function signIn(email: string, password: string) {
  supabaseLogger.info("Attempting sign in", { email });
  const { data, error } = await supabase.auth.signInWithPassword({
    email,
    password,
  });
  if (error) {
    supabaseLogger.error("Sign in failed", error);
    throw error;
  }
  supabaseLogger.info("Sign in successful", { userId: data.user?.id });
  return data;
}

/**
 * Sign up with email and password
 */
export async function signUp(email: string, password: string, fullName: string) {
  supabaseLogger.info("Attempting sign up", { email, fullName });
  const { data, error } = await supabase.auth.signUp({
    email,
    password,
    options: {
      data: {
        full_name: fullName,
      },
    },
  });
  if (error) {
    supabaseLogger.error("Sign up failed", error);
    throw error;
  }
  supabaseLogger.info("Sign up successful", { userId: data.user?.id });
  return data;
}

/**
 * Sign out
 */
export async function signOut() {
  supabaseLogger.info("Attempting sign out");
  const { error } = await supabase.auth.signOut();
  if (error) {
    supabaseLogger.error("Sign out failed", error);
    throw error;
  }
  supabaseLogger.info("Sign out successful");
}

/**
 * Listen to auth state changes
 */
export function onAuthStateChange(callback: (event: string, session: any) => void) {
  supabaseLogger.debug("Setting up auth state change listener");
  return supabase.auth.onAuthStateChange((event, session) => {
    supabaseLogger.info("Auth state changed", { event, hasSession: !!session });
    callback(event, session);
  });
}