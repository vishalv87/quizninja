"use client";

import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { ChallengeList } from "@/components/challenge/ChallengeList";
import { CreateChallengeDialog } from "@/components/challenge/CreateChallengeDialog";
import {
  usePendingChallenges,
  useActiveChallenges,
  useCompletedChallenges,
} from "@/hooks/useChallenges";
import { useChallengeStats } from "@/hooks/useChallengeStats";
import { Swords, Clock, CheckCircle2, Trophy, Target, Award } from "lucide-react";
import { Badge } from "@/components/ui/badge";

export default function ChallengesPage() {
  const [activeTab, setActiveTab] = useState("active");

  const {
    data: pendingChallenges = [],
    isLoading: pendingLoading,
    error: pendingError,
  } = usePendingChallenges();

  const {
    data: activeChallenges = [],
    isLoading: activeLoading,
    error: activeError,
  } = useActiveChallenges();

  const {
    data: completedChallenges = [],
    isLoading: completedLoading,
    error: completedError,
  } = useCompletedChallenges();

  const { data: stats } = useChallengeStats();

  return (
    <div className="container mx-auto py-8 px-4 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-start justify-between mb-4">
          <div>
            <h1 className="text-4xl font-bold mb-2">Challenges</h1>
            <p className="text-muted-foreground">
              Compete with friends and prove your knowledge!
            </p>
          </div>
          <CreateChallengeDialog />
        </div>

        {/* Stats Cards */}
        {stats && (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4 mt-6">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Challenges</CardTitle>
                <Swords className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{stats.total_challenges}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Won</CardTitle>
                <Trophy className="h-4 w-4 text-green-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-green-600">{stats.won_challenges}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Lost</CardTitle>
                <Target className="h-4 w-4 text-red-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold text-red-600">{stats.lost_challenges}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Win Rate</CardTitle>
                <Award className="h-4 w-4 text-yellow-600" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {stats.win_rate ? `${(stats.win_rate * 100).toFixed(1)}%` : "0%"}
                </div>
              </CardContent>
            </Card>
          </div>
        )}
      </div>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
        <TabsList className="grid w-full grid-cols-3 lg:w-auto lg:inline-grid">
          <TabsTrigger value="active" className="gap-2">
            <Swords className="h-4 w-4" />
            Active
            {activeChallenges.length > 0 && (
              <Badge variant="default" className="ml-1 px-2 py-0.5 text-xs">
                {activeChallenges.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="pending" className="gap-2">
            <Clock className="h-4 w-4" />
            Pending
            {pendingChallenges.length > 0 && (
              <Badge variant="secondary" className="ml-1 px-2 py-0.5 text-xs">
                {pendingChallenges.length}
              </Badge>
            )}
          </TabsTrigger>
          <TabsTrigger value="completed" className="gap-2">
            <CheckCircle2 className="h-4 w-4" />
            Completed
            {completedChallenges.length > 0 && (
              <Badge variant="outline" className="ml-1 px-2 py-0.5 text-xs">
                {completedChallenges.length}
              </Badge>
            )}
          </TabsTrigger>
        </TabsList>

        {/* Active Challenges Tab */}
        <TabsContent value="active" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Active Challenges</CardTitle>
              <CardDescription>
                Challenges that have been accepted and are in progress
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ChallengeList
                challenges={activeChallenges}
                isLoading={activeLoading}
                error={activeError}
                emptyMessage="No active challenges. Create a new challenge or wait for someone to accept yours!"
              />
            </CardContent>
          </Card>
        </TabsContent>

        {/* Pending Challenges Tab */}
        <TabsContent value="pending" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Pending Challenges</CardTitle>
              <CardDescription>
                Challenges waiting for acceptance or response
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ChallengeList
                challenges={pendingChallenges}
                isLoading={pendingLoading}
                error={pendingError}
                emptyMessage="No pending challenges. All challenges have been responded to!"
              />
            </CardContent>
          </Card>
        </TabsContent>

        {/* Completed Challenges Tab */}
        <TabsContent value="completed" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Completed Challenges</CardTitle>
              <CardDescription>
                Your challenge history and results
              </CardDescription>
            </CardHeader>
            <CardContent>
              <ChallengeList
                challenges={completedChallenges}
                isLoading={completedLoading}
                error={completedError}
                emptyMessage="No completed challenges yet. Complete some challenges to see your history!"
              />
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
