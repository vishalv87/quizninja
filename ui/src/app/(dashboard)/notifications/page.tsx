"use client";

import { useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { CheckCheck, Filter } from "lucide-react";
import { NotificationList } from "@/components/notification/NotificationList";
import {
  useNotifications,
  useNotificationStats,
  useMarkAllNotificationsAsRead,
} from "@/hooks/useNotifications";
import { getNotificationTypeLabel } from "@/lib/notification-utils";
import { NOTIFICATION_TYPES } from "@/lib/constants";

type FilterTab = "all" | "unread" | "read";
type NotificationTypeFilter = string | "all";

export default function NotificationsPage() {
  const [filterTab, setFilterTab] = useState<FilterTab>("all");
  const [typeFilter, setTypeFilter] = useState<NotificationTypeFilter>("all");

  const { data: stats } = useNotificationStats();
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
    <div className="container max-w-5xl mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Notifications</h1>
          <p className="text-muted-foreground mt-1">
            Stay updated with your latest activities
          </p>
        </div>

        {notifications.length > 0 && filterTab !== "read" && (
          <Button
            onClick={handleMarkAllAsRead}
            disabled={markAllAsRead.isPending}
            variant="outline"
          >
            <CheckCheck className="h-4 w-4 mr-2" />
            Mark all as read
          </Button>
        )}
      </div>

      {/* Stats Cards */}
      {stats && (
        <div className="grid gap-4 md:grid-cols-3">
          <Card>
            <CardHeader className="pb-3">
              <CardDescription>Total Notifications</CardDescription>
              <CardTitle className="text-3xl">
                {stats.total_notifications}
              </CardTitle>
            </CardHeader>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardDescription>Unread</CardDescription>
              <CardTitle className="text-3xl text-primary">
                {stats.unread_notifications}
              </CardTitle>
            </CardHeader>
          </Card>

          <Card>
            <CardHeader className="pb-3">
              <CardDescription>Read</CardDescription>
              <CardTitle className="text-3xl text-muted-foreground">
                {stats.read_notifications}
              </CardTitle>
            </CardHeader>
          </Card>
        </div>
      )}

      {/* Filters and Tabs */}
      <Card>
        <CardHeader>
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
            <div className="flex items-center gap-2">
              <Filter className="h-5 w-5 text-muted-foreground" />
              <CardTitle className="text-lg">Filter Notifications</CardTitle>
            </div>

            {/* Type filter */}
            <Select
              value={typeFilter}
              onValueChange={(value) => setTypeFilter(value)}
            >
              <SelectTrigger className="w-[200px]">
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
        </CardHeader>

        <CardContent>
          <Tabs
            value={filterTab}
            onValueChange={(value) => setFilterTab(value as FilterTab)}
          >
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="all">
                All
                {stats && ` (${stats.total_notifications})`}
              </TabsTrigger>
              <TabsTrigger value="unread">
                Unread
                {stats && ` (${stats.unread_notifications})`}
              </TabsTrigger>
              <TabsTrigger value="read">
                Read
                {stats && ` (${stats.read_notifications})`}
              </TabsTrigger>
            </TabsList>

            <div className="mt-6">
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
            </div>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}