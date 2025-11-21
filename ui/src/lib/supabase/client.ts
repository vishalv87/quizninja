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

// Use cookie-based client for Next.js App Router
// This ensures middleware can read the session from cookies
export const supabase = createClientComponentClient();

/**
 * Get the current session
 */
export async function getSession() {
  const { data, error } = await supabase.auth.getSession();
  if (error) {
    supabaseLogger.error("Failed to get session", error);
    throw error;
  }
  return data.session;
}

/**
 * Get the current user
 */
export async function getUser() {
  const { data, error } = await supabase.auth.getUser();
  if (error) {
    supabaseLogger.error("Failed to get user", error);
    throw error;
  }
  return data.user;
}

/**
 * Sign in with email and password
 */
export async function signIn(email: string, password: string) {
  const { data, error } = await supabase.auth.signInWithPassword({
    email,
    password,
  });
  if (error) {
    supabaseLogger.error("Sign in failed", error);
    throw error;
  }
  return data;
}

/**
 * Sign up with email and password
 */
export async function signUp(email: string, password: string, fullName: string) {
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
  return data;
}

/**
 * Sign out
 */
export async function signOut() {
  const { error } = await supabase.auth.signOut();
  if (error) {
    supabaseLogger.error("Sign out failed", error);
    throw error;
  }
}

/**
 * Listen to auth state changes
 */
export function onAuthStateChange(callback: (event: string, session: any) => void) {
  return supabase.auth.onAuthStateChange((event, session) => {
    callback(event, session);
  });
}