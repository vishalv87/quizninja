"use client";

import Link from "next/link";
import { ArrowRight, LucideIcon } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";

interface DashboardStatsCardProps {
  title: string;
  value: string | number;
  description: string;
  icon: LucideIcon;
  href: string;
  color?: string;
  isLoading?: boolean;
}

export function DashboardStatsCard({
  title,
  value,
  description,
  icon: Icon,
  href,
  color = "text-blue-600",
  isLoading = false,
}: DashboardStatsCardProps) {
  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="text-sm font-medium">{title}</CardTitle>
        <Icon className={`h-4 w-4 ${color}`} />
      </CardHeader>
      <CardContent>
        <div className="text-2xl font-bold">
          {isLoading ? "..." : value}
        </div>
        <p className="text-xs text-muted-foreground mt-1">{description}</p>
        <Button variant="ghost" size="sm" className="mt-2 h-8 px-2" asChild>
          <Link href={href}>
            View details
            <ArrowRight className="ml-1 h-3 w-3" />
          </Link>
        </Button>
      </CardContent>
    </Card>
  );
}
