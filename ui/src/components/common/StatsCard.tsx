import { cn } from "@/lib/utils";
import { LucideIcon } from "lucide-react";
import { Skeleton } from "@/components/ui/skeleton";
import { GlassCard } from "./GlassCard";

type ColorVariant = "blue" | "yellow" | "purple" | "green" | "red";

interface StatsCardProps {
  title: string;
  value: string | number;
  description?: string;
  icon: LucideIcon;
  color: ColorVariant;
  loading?: boolean;
}

const colorStyles: Record<
  ColorVariant,
  { icon: string; bg: string; border: string }
> = {
  blue: {
    icon: "text-blue-500",
    bg: "bg-blue-500/10",
    border: "border-blue-200/20",
  },
  yellow: {
    icon: "text-yellow-500",
    bg: "bg-yellow-500/10",
    border: "border-yellow-200/20",
  },
  purple: {
    icon: "text-purple-500",
    bg: "bg-purple-500/10",
    border: "border-purple-200/20",
  },
  green: {
    icon: "text-green-500",
    bg: "bg-green-500/10",
    border: "border-green-200/20",
  },
  red: {
    icon: "text-red-500",
    bg: "bg-red-500/10",
    border: "border-red-200/20",
  },
};

export function StatsCard({
  title,
  value,
  description,
  icon: Icon,
  color,
  loading = false,
}: StatsCardProps) {
  const styles = colorStyles[color];

  if (loading) {
    return (
      <GlassCard hover className={cn(styles.border)}>
        <div className="flex items-center justify-between mb-4">
          <Skeleton className="h-12 w-12 rounded-2xl" />
          <Skeleton className="h-6 w-20 rounded-full" />
        </div>
        <div className="space-y-2">
          <Skeleton className="h-8 w-16 rounded-lg" />
          <Skeleton className="h-4 w-24 rounded-lg" />
        </div>
      </GlassCard>
    );
  }

  return (
    <GlassCard hover className={cn(styles.border)}>
      <div className="flex items-center justify-between mb-4">
        <div
          className={cn(
            "p-3 rounded-2xl group-hover:scale-110 transition-transform duration-300 ring-1 ring-white/20",
            styles.bg
          )}
        >
          <Icon className={cn("h-6 w-6", styles.icon)} />
        </div>
        <span className="text-xs font-bold text-slate-500 dark:text-slate-400 bg-white/50 dark:bg-white/10 px-2.5 py-1 rounded-full border border-white/20 backdrop-blur-sm">
          {title}
        </span>
      </div>
      <div className="space-y-1">
        <h3 className="text-3xl font-bold tracking-tight text-slate-800 dark:text-slate-100">
          {value}
        </h3>
        {description && (
          <p className="text-sm text-slate-500 dark:text-slate-400 font-medium">
            {description}
          </p>
        )}
      </div>
    </GlassCard>
  );
}