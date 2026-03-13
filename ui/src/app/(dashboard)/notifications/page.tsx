"use client";

import { useState } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Button } from "@/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CheckCheck, Filter, Bell, Mail, MailOpen } from "lucide-react";
import { NotificationList } from "@/components/notification/NotificationList";
import {
  useNotifications,
  useNotificationStats,
  useMarkAllNotificationsAsRead,
} from "@/hooks/useNotifications";
import { getNotificationTypeLabel } from "@/lib/notification-utils";
import { NOTIFICATION_TYPES } from "@/lib/constants";
import { type NotificationType } from "@/constants";
import { GlassCard } from "@/components/common/GlassCard";
import { StatsCard } from "@/components/common/StatsCard";
import { StatsGrid } from "@/components/common/StatsGrid";

type FilterTab = "all" | "unread" | "read";
type NotificationTypeFilter = NotificationType | "all";

export default function NotificationsPage() {
  const [filterTab, setFilterTab] = useState<FilterTab>("all");
  const [typeFilter, setTypeFilter] = useState<NotificationTypeFilter>("all");

  const { data: stats, isLoading: statsLoading } = useNotificationStats();
  const markAllAsRead = useMarkAllNotificationsAsRead();

  // Build filter based on selected tab and type
  const filters =
    filterTab === "all"
      ? { type: typeFilter !== "all" ? typeFilter : undefined }
      : {
          is_read: filterTab === "read",
          type: typeFilter !== "all" ? typeFilter : undefined,
        };

  const { data: notifications = [], isLoading } = useNotifications(filters);

  const handleMarkAllAsRead = () => {
    markAllAsRead.mutate();
  };

  return (
    <div className="space-y-10 pb-10">
      {/* Header with action button */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
            Notifications
          </h1>
          <p className="text-slate-500 dark:text-slate-400 mt-1">
            Stay updated with your latest activities and achievements
          </p>
        </div>
        {notifications.length > 0 && filterTab !== "read" && (
          <Button
            onClick={handleMarkAllAsRead}
            disabled={markAllAsRead.isPending}
            size="lg"
            className="bg-gradient-to-r from-violet-600 to-indigo-600 hover:from-violet-700 hover:to-indigo-700 text-white font-semibold h-11 px-6 rounded-xl shadow-lg shadow-indigo-500/25 transition-all hover:shadow-xl"
          >
            <CheckCheck className="h-4 w-4 mr-2" />
            Mark all as read
          </Button>
        )}
      </div>

      {/* Stats Cards */}
      <StatsGrid columns={3}>
        <StatsCard
          title="Total"
          value={stats?.total_notifications ?? 0}
          description="All notifications"
          icon={Bell}
          color="blue"
          loading={statsLoading}
        />
        <StatsCard
          title="Unread"
          value={stats?.unread_notifications ?? 0}
          description="Awaiting your attention"
          icon={Mail}
          color="yellow"
          loading={statsLoading}
        />
        <StatsCard
          title="Read"
          value={stats?.read_notifications ?? 0}
          description="Already reviewed"
          icon={MailOpen}
          color="green"
          loading={statsLoading}
        />
      </StatsGrid>

      {/* Filters and Content */}
      <div className="container px-0 md:px-4">
        <GlassCard padding="none" rounded="2xl">
          {/* Filter Header */}
          <div className="p-6 border-b border-white/10">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div className="flex items-center gap-2">
                <span className="bg-gradient-to-br from-violet-500 to-purple-600 text-white p-1.5 rounded-lg shadow-sm">
                  <Filter className="h-4 w-4" />
                </span>
                <h2 className="text-lg font-bold tracking-tight text-slate-800 dark:text-slate-100">
                  Filter Notifications
                </h2>
              </div>

              {/* Type filter */}
              <Select
                value={typeFilter}
                onValueChange={(value) => setTypeFilter(value as NotificationTypeFilter)}
              >
                <SelectTrigger className="w-[200px] bg-white/50 dark:bg-white/10 border-white/20 dark:border-white/10 backdrop-blur-sm">
                  <SelectValue placeholder="All types" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All types</SelectItem>
                  {Object.values(NOTIFICATION_TYPES).map((type) => (
                    <SelectItem key={type} value={type}>
                      {getNotificationTypeLabel(type)}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
          </div>

          {/* Tabs and Content */}
          <div className="p-6">
            <Tabs
              value={filterTab}
              onValueChange={(value) => setFilterTab(value as FilterTab)}
            >
              <TabsList className="grid w-full max-w-md grid-cols-3 bg-white/60 dark:bg-black/40 backdrop-blur-md border border-white/20 dark:border-white/10 p-1 rounded-xl shadow-sm mb-6">
                <TabsTrigger
                  value="all"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  All
                  {stats && ` (${stats.total_notifications})`}
                </TabsTrigger>
                <TabsTrigger
                  value="unread"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  Unread
                  {stats && ` (${stats.unread_notifications})`}
                </TabsTrigger>
                <TabsTrigger
                  value="read"
                  className="rounded-lg data-[state=active]:bg-white/90 dark:data-[state=active]:bg-background/90 data-[state=active]:text-violet-700 dark:data-[state=active]:text-violet-400 data-[state=active]:shadow-md data-[state=active]:border data-[state=active]:border-violet-200/50 dark:data-[state=active]:border-violet-800/50 transition-all duration-300 hover:bg-white/40 dark:hover:bg-white/5"
                >
                  Read
                  {stats && ` (${stats.read_notifications})`}
                </TabsTrigger>
              </TabsList>

              <TabsContent value="all" className="mt-0">
                <NotificationList
                  notifications={notifications}
                  isLoading={isLoading}
                  emptyMessage="You don't have any notifications yet. When you receive notifications, they'll appear here."
                  groupByDate
                />
              </TabsContent>

              <TabsContent value="unread" className="mt-0">
                <NotificationList
                  notifications={notifications}
                  isLoading={isLoading}
                  emptyMessage="You're all caught up! No unread notifications."
                  groupByDate
                />
              </TabsContent>

              <TabsContent value="read" className="mt-0">
                <NotificationList
                  notifications={notifications}
                  isLoading={isLoading}
                  emptyMessage="No read notifications yet."
                  groupByDate
                />
              </TabsContent>
            </Tabs>
          </div>
        </GlassCard>
      </div>
    </div>
  );
}
