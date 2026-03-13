import { cn } from "@/lib/utils";

interface StatsGridProps {
  children: React.ReactNode;
  columns?: 2 | 3 | 4;
  className?: string;
}

const columnStyles = {
  2: "md:grid-cols-2",
  3: "md:grid-cols-2 lg:grid-cols-3",
  4: "md:grid-cols-2 lg:grid-cols-4",
};

export function StatsGrid({ children, columns = 4, className }: StatsGridProps) {
  return (
    <div className={cn("grid gap-6", columnStyles[columns], className)}>
      {children}
    </div>
  );
}