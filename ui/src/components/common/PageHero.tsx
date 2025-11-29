import { cn } from "@/lib/utils";

type GradientVariant = "default" | "blue" | "green" | "amber";

interface PageHeroProps {
  title: string;
  description: string;
  icon?: React.ReactNode;
  children?: React.ReactNode;
  gradient?: GradientVariant;
}

const gradientStyles: Record<GradientVariant, string> = {
  default: "from-violet-600 via-indigo-600 to-purple-700 shadow-indigo-500/30",
  blue: "from-blue-600 via-cyan-600 to-blue-700 shadow-blue-500/30",
  green: "from-emerald-600 via-teal-600 to-green-700 shadow-emerald-500/30",
  amber: "from-amber-500 via-orange-500 to-red-600 shadow-orange-500/30",
};

export function PageHero({
  title,
  description,
  icon,
  children,
  gradient = "default",
}: PageHeroProps) {
  return (
    <div
      className={cn(
        "relative overflow-hidden rounded-3xl p-8 text-white shadow-2xl lg:p-12 border border-white/10",
        "bg-gradient-to-br",
        gradientStyles[gradient]
      )}
    >
      {/* Content */}
      <div className="relative z-10 max-w-2xl">
        <h1 className="text-4xl font-bold tracking-tight sm:text-5xl mb-6 drop-shadow-sm flex items-center gap-3">
          {icon && <span>{icon}</span>}
          {title}
        </h1>
        <p className="text-xl text-white/90 mb-8 font-medium leading-relaxed">
          {description}
        </p>
        {children}
      </div>

      {/* Decorative background elements */}
      <div className="absolute right-0 top-0 -mt-20 -mr-20 h-96 w-96 rounded-full bg-white/10 blur-3xl" />
      <div className="absolute bottom-0 right-20 -mb-20 h-64 w-64 rounded-full bg-white/20 blur-3xl" />
      <div className="absolute left-10 bottom-10 h-32 w-32 rounded-full bg-white/10 blur-2xl" />
    </div>
  );
}