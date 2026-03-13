import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api/client";
import { API_ENDPOINTS } from "@/lib/api/endpoints";
import { toast } from "sonner";

interface OnboardingStatus {
  onboarding_completed: boolean;
}

/**
 * Hook to get onboarding status
 */
export function useOnboardingStatus() {
  return useQuery({
    queryKey: ["onboarding-status"],
    queryFn: async () => {
      const response = await apiClient.get<OnboardingStatus>(
        API_ENDPOINTS.USERS.ONBOARDING.STATUS
      ) as unknown as OnboardingStatus;
      return response;
    },
    staleTime: Infinity, // Cache forever until invalidated
  });
}

/**
 * Hook to mark onboarding as complete
 */
export function useCompleteOnboarding() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await apiClient.post(API_ENDPOINTS.USERS.ONBOARDING.COMPLETE);
    },
    onSuccess: () => {
      // Invalidate onboarding status
      queryClient.invalidateQueries({ queryKey: ["onboarding-status"] });
      // Invalidate user profile as it includes onboarding status
      queryClient.invalidateQueries({ queryKey: ["profile"] });
      queryClient.invalidateQueries({ queryKey: ["preferences"] });
      toast.success("Welcome to QuizNinja!", {
        description: "Your account setup is complete.",
      });
    },
    onError: (error: any) => {
      toast.error("Failed to complete onboarding", {
        description: error.message || "Could not complete onboarding.",
      });
    },
  });
}
