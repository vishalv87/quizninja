import { useEffect, useRef } from "react";
import { useQuizStore } from "@/store/quizStore";

/**
 * Hook to manage quiz timer
 * Automatically decrements time every second
 *
 * @param onTimeUp - Callback when time runs out
 */
export function useQuizTimer(onTimeUp?: () => void) {
  const { timeRemaining, decrementTime } = useQuizStore();
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    // Don't start timer if no time limit
    if (timeRemaining === null) {
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
  }, [timeRemaining, decrementTime, onTimeUp]);

  return {
    timeRemaining,
  };
}
