"use client";

import { useState } from "react";
import Link from "next/link";
import Image from "next/image";
import { BookOpen } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import type { Category } from "@/types/quiz";

interface CategoryCardProps {
  category: Category;
}

export function CategoryCard({ category }: CategoryCardProps) {
  const [imageError, setImageError] = useState(false);

  return (
    <Link href={`/quizzes/category/${category.id}`}>
      <Card className="hover:shadow-lg transition-all hover:scale-105 cursor-pointer h-full">
        <CardHeader>
          <div className="flex items-center justify-between mb-2">
            {category.icon_url && !imageError ? (
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center overflow-hidden relative">
                <Image
                  src={category.icon_url}
                  alt={category.display_name}
                  fill
                  className="object-cover"
                  onError={() => setImageError(true)}
                  unoptimized
                />
              </div>
            ) : (
              <div className="w-12 h-12 rounded-lg bg-primary/10 flex items-center justify-center">
                <BookOpen className="h-6 w-6 text-primary" />
              </div>
            )}
            <Badge variant="secondary" className="font-semibold">
              {category.quiz_count} {category.quiz_count === 1 ? "Quiz" : "Quizzes"}
            </Badge>
          </div>
          <CardTitle className="text-xl">{category.display_name}</CardTitle>
          <CardDescription className="line-clamp-2 min-h-[2.5rem]">
            {category.description || "Explore quizzes in this category"}
          </CardDescription>
        </CardHeader>
      </Card>
    </Link>
  );
}
