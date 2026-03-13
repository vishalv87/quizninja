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
  const hasCalledOnTimeUp = useRef(false);

  useEffect(() => {
    // Don't start timer if no time limit
    if (timeRemaining === null) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      return;
    }

    // Check if time is up - only call onTimeUp once
    if (timeRemaining <= 0) {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
        intervalRef.current = null;
      }
      // Guard to prevent multiple calls to onTimeUp
      if (!hasCalledOnTimeUp.current) {
        hasCalledOnTimeUp.current = true;
        onTimeUp?.();
      }
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
