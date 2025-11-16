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
  Target,
  Zap,
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

  const platformStats = [
    {
      title: "Quiz Categories",
      value: "50+",
      description: "Topics to explore",
    },
    {
      title: "Active Players",
      value: "10K+",
      description: "Join the community",
    },
    {
      title: "Achievements",
      value: "100+",
      description: "Badges to unlock",
    },
    {
      title: "Daily Quizzes",
      value: "500+",
      description: "New challenges daily",
    },
  ];

  const features = [
    {
      icon: Target,
      title: "Take Quizzes",
      description: "Test your knowledge across various categories and difficulty levels",
      iconBg: "bg-blue-100",
      iconColor: "text-blue-600",
    },
    {
      icon: Users,
      title: "Challenge Friends",
      description: "Compete with friends and climb the leaderboard together",
      iconBg: "bg-green-100",
      iconColor: "text-green-600",
    },
    {
      icon: Trophy,
      title: "Earn Achievements",
      description: "Unlock badges and achievements as you progress",
      iconBg: "bg-yellow-100",
      iconColor: "text-yellow-600",
    },
    {
      icon: MessageSquare,
      title: "Join Discussions",
      description: "Engage with the community and share your insights",
      iconBg: "bg-purple-100",
      iconColor: "text-purple-600",
    },
  ];

  return (
    <div className="min-h-screen bg-background">
      {/* Header Logo */}
      <div className="mx-auto px-6 py-6 max-w-7xl w-full">
        <div className="text-2xl font-bold">QuizNinja</div>
      </div>

      <div className="mx-auto px-6 lg:px-8 py-12 max-w-7xl w-full">
        <div className="space-y-16">
          {/* Hero Section */}
          <div className="text-center space-y-6">
            <div className="inline-flex items-center justify-center w-20 h-20 bg-gray-900 rounded-2xl mb-4">
              <Target className="w-10 h-10 text-white" />
            </div>
            <h1 className="text-6xl font-bold tracking-tight">
              Welcome to <span className="bg-clip-text text-transparent bg-gradient-to-r from-blue-600 to-purple-600">QuizNinja</span>
            </h1>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              Your journey to knowledge mastery starts here. Challenge yourself, compete with friends, and unlock achievements!
            </p>
          </div>

          {/* Platform Stats */}
          <div className="grid gap-6 grid-cols-2 lg:grid-cols-4">
            {platformStats.map((stat, index) => (
              <Card key={index} className="bg-white border border-gray-200 hover:shadow-md transition-shadow">
                <CardContent className="p-6 text-center">
                  <div className="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2">
                    {stat.title}
                  </div>
                  <div className="text-4xl font-bold mb-1">{stat.value}</div>
                  <p className="text-sm text-gray-600">
                    {stat.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>

          {/* Features Section */}
          <div className="space-y-8">
            <div className="text-center">
              <h2 className="text-3xl font-bold tracking-tight">
                Everything You Need to Excel
              </h2>
              <p className="text-gray-600 mt-2">
                Powerful features to enhance your learning experience
              </p>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
              {features.map((feature, index) => {
                const Icon = feature.icon;
                return (
                  <Card key={index} className="bg-white border border-gray-200 hover:shadow-md transition-shadow">
                    <CardContent className="p-6">
                      <div className="flex flex-col items-start gap-4">
                        <div className={`flex items-center justify-center w-12 h-12 rounded-full ${feature.iconBg}`}>
                          <Icon className={`w-6 h-6 ${feature.iconColor}`} />
                        </div>
                        <div>
                          <h3 className="font-bold text-lg mb-2">{feature.title}</h3>
                          <p className="text-sm text-gray-600">
                            {feature.description}
                          </p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </div>
          </div>

          {/* CTA Section */}
          <div className="space-y-4 pt-4">
            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <Button
                size="lg"
                className="w-full sm:w-auto bg-gray-900 hover:bg-gray-800 text-white h-14 px-8 rounded-lg"
                onClick={handleGetStarted}
              >
                Get Started
                <ArrowRight className="ml-2 h-5 w-5" />
              </Button>
              <Button
                size="lg"
                variant="outline"
                className="w-full sm:w-auto bg-gray-100 hover:bg-gray-200 text-gray-900 border-0 h-14 px-8 rounded-lg"
                onClick={handleSkip}
                disabled={completeOnboarding.isPending}
              >
                {completeOnboarding.isPending ? "Skipping..." : "Skip for Now"}
              </Button>
            </div>

            <p className="text-center text-sm text-gray-600">
              Set up your preferences to get personalized quiz recommendations
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
