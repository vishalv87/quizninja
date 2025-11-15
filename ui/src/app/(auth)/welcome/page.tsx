"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  BookOpen,
  Users,
  Trophy,
  MessageSquare,
  ArrowRight,
  Sparkles,
} from "lucide-react";
import { useOnboardingStatus, useCompleteOnboarding } from "@/hooks/useOnboarding";

export default function WelcomePage() {
  const router = useRouter();
  const { data: onboardingStatus } = useOnboardingStatus();
  const completeOnboarding = useCompleteOnboarding();

  // Redirect to dashboard if onboarding is already complete
  useEffect(() => {
    if (onboardingStatus?.onboarding_completed) {
      router.push("/dashboard");
    }
  }, [onboardingStatus, router]);

  const handleGetStarted = () => {
    router.push("/preferences");
  };

  const handleSkip = () => {
    completeOnboarding.mutate(undefined, {
      onSuccess: () => {
        router.push("/dashboard");
      },
    });
  };

  const features = [
    {
      icon: BookOpen,
      title: "Take Quizzes",
      description: "Test your knowledge across various categories and difficulty levels",
    },
    {
      icon: Users,
      title: "Challenge Friends",
      description: "Compete with friends and climb the leaderboard together",
    },
    {
      icon: Trophy,
      title: "Earn Achievements",
      description: "Unlock badges and achievements as you progress",
    },
    {
      icon: MessageSquare,
      title: "Join Discussions",
      description: "Engage with the community and share your insights",
    },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary/5 via-background to-secondary/5 flex items-center justify-center p-4">
      <div className="max-w-4xl w-full space-y-8">
        {/* Hero Section */}
        <div className="text-center space-y-4">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-primary rounded-2xl mb-4">
            <Sparkles className="w-10 h-10 text-primary-foreground" />
          </div>
          <h1 className="text-5xl md:text-6xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-primary to-primary/60">
            Welcome to QuizNinja
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Your journey to knowledge mastery starts here. Challenge yourself, compete with friends, and unlock achievements!
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {features.map((feature, index) => {
            const Icon = feature.icon;
            return (
              <Card key={index} className="border-2 hover:border-primary/50 transition-colors">
                <CardContent className="p-6">
                  <div className="flex items-start gap-4">
                    <div className="flex items-center justify-center w-12 h-12 bg-primary/10 rounded-lg flex-shrink-0">
                      <Icon className="w-6 h-6 text-primary" />
                    </div>
                    <div>
                      <h3 className="font-semibold text-lg mb-1">{feature.title}</h3>
                      <p className="text-sm text-muted-foreground">{feature.description}</p>
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })}
        </div>

        {/* CTA Buttons */}
        <div className="flex flex-col sm:flex-row gap-4 justify-center items-center pt-8">
          <Button
            size="lg"
            className="w-full sm:w-auto text-lg h-14 px-8"
            onClick={handleGetStarted}
          >
            Get Started
            <ArrowRight className="ml-2 h-5 w-5" />
          </Button>
          <Button
            size="lg"
            variant="outline"
            className="w-full sm:w-auto text-lg h-14 px-8"
            onClick={handleSkip}
            disabled={completeOnboarding.isPending}
          >
            {completeOnboarding.isPending ? "Skipping..." : "Skip for Now"}
          </Button>
        </div>

        {/* Helper Text */}
        <p className="text-center text-sm text-muted-foreground">
          Set up your preferences to get personalized quiz recommendations
        </p>
      </div>
    </div>
  );
}
