import { cn } from "@/lib/utils";

interface GlassCardProps {
  children: React.ReactNode;
  hover?: boolean;
  padding?: "none" | "sm" | "md" | "lg";
  rounded?: "xl" | "2xl" | "3xl";
  className?: string;
}

const paddingStyles = {
  none: "",
  sm: "p-4",
  md: "p-6",
  lg: "p-8",
};

const roundedStyles = {
  xl: "rounded-xl",
  "2xl": "rounded-2xl",
  "3xl": "rounded-3xl",
};

export function GlassCard({
  children,
  hover = false,
  padding = "md",
  rounded = "2xl",
  className,
}: GlassCardProps) {
  return (
    <div
      className={cn(
        // Base glassmorphism styling
        "bg-white/40 dark:bg-black/40",
        "backdrop-blur-md",
        "border border-white/20 dark:border-white/10",
        "shadow-lg shadow-black/5",
        "transition-all duration-300",
        // Padding
        paddingStyles[padding],
        // Rounded corners
        roundedStyles[rounded],
        // Hover effects (optional)
        hover && "hover:shadow-xl hover:-translate-y-1",
        // Custom className
        className
      )}
    >
      {children}
    </div>
  );
}
