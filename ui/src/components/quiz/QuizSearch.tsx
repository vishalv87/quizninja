"use client";

import { Input } from "@/components/ui/input";
import { Search, X } from "lucide-react";
import { Button } from "@/components/ui/button";

import { cn } from "@/lib/utils";

interface QuizSearchProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export function QuizSearch({
  value,
  onChange,
  placeholder = "Search quizzes...",
  className,
}: QuizSearchProps) {
  const handleClear = () => {
    onChange("");
  };

  return (
    <div className="relative w-full">
      <Search className={cn("absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground", className && "text-inherit opacity-70")} />
      <Input
        type="text"
        placeholder={placeholder}
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className={cn("pl-10 pr-10", className)}
      />
      {value && (
        <Button
          variant="ghost"
          size="icon"
          className={cn("absolute right-1 top-1/2 -translate-y-1/2 h-7 w-7", className && "text-inherit opacity-70 hover:bg-white/20 hover:text-inherit")}
          onClick={handleClear}
        >
          <X className="h-4 w-4" />
          <span className="sr-only">Clear search</span>
        </Button>
      )}
    </div>
  );
}