"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ArrowLeft, ArrowRight, Target, Check } from "lucide-react";
import { useCategories } from "@/hooks/useCategories";
import { useUpdatePreferences } from "@/hooks/usePreferences";
import { useCompleteOnboarding } from "@/hooks/useOnboarding";
import { Skeleton } from "@/components/ui/skeleton";
import { toast } from "sonner";

export default function PreferencesPage() {
  const router = useRouter();
  const [selectedCategories, setSelectedCategories] = useState<string[]>([]);
  const [difficulty, setDifficulty] = useState<string>("medium");

  const { data: categories, isLoading: categoriesLoading } = useCategories();
  const updatePreferences = useUpdatePreferences();
  const completeOnboarding = useCompleteOnboarding();

  const handleCategoryToggle = (categoryId: string) => {
    setSelectedCategories((prev) =>
      prev.includes(categoryId)
        ? prev.filter((id) => id !== categoryId)
        : [...prev, categoryId]
    );
  };

  const handleBack = () => {
    router.push("/");
  };

  const handleComplete = async () => {
    // Validate selection
    if (selectedCategories.length === 0) {
      toast.error("Please select at least one category");
      return;
    }

    // Update preferences first
    try {
      await updatePreferences.mutateAsync({
        category_preferences: selectedCategories,
        difficulty_level: difficulty,
      });

      // Then mark onboarding as complete
      completeOnboarding.mutate(undefined, {
        onSuccess: () => {
          router.push("/dashboard");
        },
      });
    } catch (error) {
      console.error("Failed to save preferences:", error);
    }
  };

  const isLoading = updatePreferences.isPending || completeOnboarding.isPending;

  return (
    <div className="min-h-screen bg-background">
      <div className="mx-auto px-6 lg:px-8 py-12 max-w-7xl w-full">
        <div className="space-y-16">
          {/* Header */}
          <div className="text-center space-y-6">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-gray-900 rounded-2xl mb-4">
              <Target className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-6xl font-bold tracking-tight">
              Customize Your <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">Experience</span>
            </h1>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              Select your preferred categories and difficulty level to get personalized quiz recommendations
            </p>
          </div>

          {/* Categories Selection */}
          <div className="space-y-8">
            <div className="text-center">
              <h2 className="text-3xl font-bold tracking-tight">Favorite Categories</h2>
              <p className="text-gray-600 mt-2">
                Choose the topics you're most interested in (select at least one)
              </p>
            </div>

            {categoriesLoading ? (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {Array.from({ length: 6 }).map((_, i) => (
                  <Skeleton key={i} className="h-32 w-full" />
                ))}
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
                {categories?.map((category, index) => {
                  const colors = [
                    { bg: "bg-blue-100", text: "text-blue-600", border: "border-blue-600" },
                    { bg: "bg-green-100", text: "text-green-600", border: "border-green-600" },
                    { bg: "bg-yellow-100", text: "text-yellow-600", border: "border-yellow-600" },
                    { bg: "bg-purple-100", text: "text-purple-600", border: "border-purple-600" },
                    { bg: "bg-pink-100", text: "text-pink-600", border: "border-pink-600" },
                    { bg: "bg-indigo-100", text: "text-indigo-600", border: "border-indigo-600" },
                  ];
                  const colorScheme = colors[index % colors.length];
                  const isSelected = selectedCategories.includes(category.id);

                  return (
                    <Card
                      key={category.id}
                      className={`cursor-pointer hover:shadow-md transition-all ${
                        isSelected
                          ? `bg-white border-2 ${colorScheme.border}`
                          : "bg-white border border-gray-200"
                      }`}
                      onClick={() => handleCategoryToggle(category.id)}
                    >
                      <CardContent className="p-6">
                        <div className="flex flex-col items-center gap-4 text-center">
                          <div className={`flex items-center justify-center w-16 h-16 rounded-full ${colorScheme.bg}`}>
                            {isSelected ? (
                              <Check className={`w-8 h-8 ${colorScheme.text}`} />
                            ) : (
                              <span className={`text-2xl font-bold ${colorScheme.text}`}>
                                {category.display_name.charAt(0)}
                              </span>
                            )}
                          </div>
                          <div>
                            <h3 className="font-bold text-lg">{category.display_name}</h3>
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  );
                })}
              </div>
            )}

            {selectedCategories.length > 0 && (
              <p className="text-center text-sm text-gray-600">
                {selectedCategories.length} {selectedCategories.length === 1 ? "category" : "categories"} selected
              </p>
            )}
          </div>

          {/* Difficulty Level */}
          <div className="space-y-8">
            <div className="text-center">
              <h2 className="text-3xl font-bold tracking-tight">Preferred Difficulty</h2>
              <p className="text-gray-600 mt-2">
                Choose your preferred challenge level
              </p>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
              {[
                {
                  value: "easy",
                  label: "Easy",
                  description: "Perfect for beginners and casual learning",
                  color: { bg: "bg-green-100", text: "text-green-600", border: "border-green-600" },
                },
                {
                  value: "medium",
                  label: "Medium",
                  description: "Balanced challenge for most learners",
                  color: { bg: "bg-yellow-100", text: "text-yellow-600", border: "border-yellow-600" },
                },
                {
                  value: "hard",
                  label: "Hard",
                  description: "For experts seeking maximum challenge",
                  color: { bg: "bg-red-100", text: "text-red-600", border: "border-red-600" },
                },
              ].map((level) => {
                const isSelected = difficulty === level.value;
                return (
                  <Card
                    key={level.value}
                    className={`cursor-pointer hover:shadow-md transition-all ${
                      isSelected
                        ? `bg-white border-2 ${level.color.border}`
                        : "bg-white border border-gray-200"
                    }`}
                    onClick={() => setDifficulty(level.value)}
                  >
                    <CardContent className="p-6 text-center">
                      <div className="flex flex-col items-center gap-4">
                        <div className={`flex items-center justify-center w-16 h-16 rounded-full ${level.color.bg}`}>
                          {isSelected ? (
                            <Check className={`w-8 h-8 ${level.color.text}`} />
                          ) : (
                            <span className={`text-2xl font-bold ${level.color.text}`}>
                              {level.label.charAt(0)}
                            </span>
                          )}
                        </div>
                        <div>
                          <h3 className="font-bold text-xl mb-2">{level.label}</h3>
                          <p className="text-sm text-gray-600">
                            {level.description}
                          </p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          </div>

          {/* Action Buttons */}
          <div className="space-y-4 pt-4">
            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <Button
                size="lg"
                variant="outline"
                className="w-full sm:w-auto bg-gray-100 hover:bg-gray-200 text-gray-900 border-0 h-14 px-8 rounded-lg"
                onClick={handleBack}
                disabled={isLoading}
              >
                <ArrowLeft className="mr-2 h-5 w-5" />
                Back
              </Button>
              <Button
                size="lg"
                className="w-full sm:w-auto bg-gray-900 hover:bg-gray-800 text-white h-14 px-8 rounded-lg"
                onClick={handleComplete}
                disabled={isLoading || selectedCategories.length === 0}
              >
                {isLoading ? "Saving..." : "Complete Setup"}
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
            </div>

            {/* Helper Text */}
            <p className="text-center text-sm text-gray-600">
              You can always change these preferences later in your settings
            </p>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="mx-auto px-6 py-8 max-w-7xl w-full">
        <p className="text-sm text-gray-500">
          © 2024 QuizNinja. All rights reserved.
        </p>
      </div>
    </div>
  );
}
