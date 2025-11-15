"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { ArrowLeft, ArrowRight, Sparkles } from "lucide-react";
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
    router.push("/welcome");
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
    <div className="min-h-screen bg-gradient-to-br from-primary/5 via-background to-secondary/5 flex items-center justify-center p-4">
      <div className="max-w-3xl w-full space-y-6">
        {/* Header */}
        <div className="text-center space-y-2">
          <div className="inline-flex items-center justify-center w-16 h-16 bg-primary rounded-2xl mb-3">
            <Sparkles className="w-8 h-8 text-primary-foreground" />
          </div>
          <h1 className="text-4xl font-bold">Customize Your Experience</h1>
          <p className="text-muted-foreground">
            Select your preferred categories and difficulty level to get personalized quiz recommendations
          </p>
        </div>

        {/* Categories Selection */}
        <Card>
          <CardHeader>
            <CardTitle>Favorite Categories</CardTitle>
            <CardDescription>
              Choose the topics you're most interested in (select at least one)
            </CardDescription>
          </CardHeader>
          <CardContent>
            {categoriesLoading ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {Array.from({ length: 6 }).map((_, i) => (
                  <Skeleton key={i} className="h-12 w-full" />
                ))}
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {categories?.map((category) => (
                  <div
                    key={category.id}
                    className="flex items-center space-x-3 p-3 rounded-lg border hover:bg-accent transition-colors cursor-pointer"
                    onClick={() => handleCategoryToggle(category.id)}
                  >
                    <Checkbox
                      id={category.id}
                      checked={selectedCategories.includes(category.id)}
                      onCheckedChange={() => handleCategoryToggle(category.id)}
                    />
                    <Label
                      htmlFor={category.id}
                      className="flex-1 cursor-pointer font-medium"
                    >
                      {category.display_name}
                    </Label>
                  </div>
                ))}
              </div>
            )}
            {selectedCategories.length > 0 && (
              <p className="text-sm text-muted-foreground mt-4">
                {selectedCategories.length} {selectedCategories.length === 1 ? "category" : "categories"} selected
              </p>
            )}
          </CardContent>
        </Card>

        {/* Difficulty Level */}
        <Card>
          <CardHeader>
            <CardTitle>Preferred Difficulty</CardTitle>
            <CardDescription>
              Choose your preferred challenge level
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Select value={difficulty} onValueChange={setDifficulty}>
              <SelectTrigger className="w-full">
                <SelectValue placeholder="Select difficulty level" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="easy">
                  <div className="flex flex-col items-start">
                    <span className="font-medium">Easy</span>
                    <span className="text-xs text-muted-foreground">
                      Perfect for beginners and casual learning
                    </span>
                  </div>
                </SelectItem>
                <SelectItem value="medium">
                  <div className="flex flex-col items-start">
                    <span className="font-medium">Medium</span>
                    <span className="text-xs text-muted-foreground">
                      Balanced challenge for most learners
                    </span>
                  </div>
                </SelectItem>
                <SelectItem value="hard">
                  <div className="flex flex-col items-start">
                    <span className="font-medium">Hard</span>
                    <span className="text-xs text-muted-foreground">
                      For experts seeking maximum challenge
                    </span>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </CardContent>
        </Card>

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-between pt-4">
          <Button
            size="lg"
            variant="outline"
            onClick={handleBack}
            disabled={isLoading}
          >
            <ArrowLeft className="mr-2 h-5 w-5" />
            Back
          </Button>
          <Button
            size="lg"
            onClick={handleComplete}
            disabled={isLoading || selectedCategories.length === 0}
            className="sm:w-auto"
          >
            {isLoading ? "Saving..." : "Complete Setup"}
            <ArrowRight className="ml-2 h-5 w-5" />
          </Button>
        </div>

        {/* Helper Text */}
        <p className="text-center text-sm text-muted-foreground">
          You can always change these preferences later in your settings
        </p>
      </div>
    </div>
  );
}
