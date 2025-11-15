import { useEffect, useRef } from "react";
import { useQuizStore } from "@/store/quizStore";

/**
 * Hook to manage quiz timer
 * Automatically decrements time every second when not paused
 *
 * @param onTimeUp - Callback when time runs out
 */
export function useQuizTimer(onTimeUp?: () => void) {
  const { timeRemaining, isPaused, decrementTime } = useQuizStore();
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Don't start timer if no time limit or paused
    if (timeRemaining === null || isPaused) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      return;
    }

    // Check if time is up
    if (timeRemaining <= 0) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      onTimeUp?.();
      return;
    }

    // Start countdown
    intervalRef.current = setInterval(() => {
      decrementTime();
    }, 1000);

    // Cleanup on unmount
    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [timeRemaining, isPaused, decrementTime, onTimeUp]);

  return {
    timeRemaining,
    isPaused,
  };
}
