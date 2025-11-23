"use client";

import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import type { Question } from "@/types/quiz";
import { cn } from "@/lib/utils";

interface QuestionCardProps {
  question: Question;
  questionNumber: number;
  selectedAnswer?: string;
  onAnswerChange: (answer: string) => void;
  disabled?: boolean;
}

export function QuestionCard({
  question,
  questionNumber,
  selectedAnswer,
  onAnswerChange,
  disabled = false,
}: QuestionCardProps) {
  // Normalize question type to handle both camelCase and snake_case from backend
  const questionType = question.question_type?.toLowerCase().replace(/_/g, '');
  const isMultipleChoice = questionType === "multiplechoice";
  const isTrueFalse = questionType === "truefalse";
  const isShortAnswer = questionType === "shortanswer";

  // Render multiple choice question
  if (isMultipleChoice && question.options) {
    return (
      <Card>
        <CardHeader>
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center font-bold">
              {questionNumber}
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-semibold leading-relaxed">
                {question.question_text}
              </h3>
              <p className="text-sm text-muted-foreground mt-1">
                {question.points} {question.points === 1 ? "point" : "points"}
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <RadioGroup
            value={selectedAnswer ?? ""}
            onValueChange={onAnswerChange}
            disabled={disabled}
            className="space-y-3"
          >
            {question.options.map((optionText, index) => {
              const optionId = `${question.id}-option-${index}`;
              return (
                <div
                  key={optionId}
                  className={cn(
                    "flex items-center space-x-3 p-4 rounded-lg border-2 transition-all",
                    selectedAnswer === optionText
                      ? "border-primary bg-primary/5"
                      : "border-border hover:border-primary/50"
                  )}
                >
                  <RadioGroupItem value={optionText} id={optionId} />
                  <Label
                    htmlFor={optionId}
                    className="flex-1 cursor-pointer text-base"
                  >
                    {optionText}
                  </Label>
                </div>
              );
            })}
          </RadioGroup>
        </CardContent>
      </Card>
    );
  }

  // Render true/false question
  if (isTrueFalse) {
    return (
      <Card>
        <CardHeader>
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center font-bold">
              {questionNumber}
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-semibold leading-relaxed">
                {question.question_text}
              </h3>
              <p className="text-sm text-muted-foreground mt-1">
                {question.points} {question.points === 1 ? "point" : "points"}
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <RadioGroup
            value={selectedAnswer ?? ""}
            onValueChange={onAnswerChange}
            disabled={disabled}
            className="space-y-3"
          >
            <div
              className={cn(
                "flex items-center space-x-3 p-4 rounded-lg border-2 transition-all",
                selectedAnswer === "true"
                  ? "border-primary bg-primary/5"
                  : "border-border hover:border-primary/50"
              )}
            >
              <RadioGroupItem value="true" id="true" />
              <Label htmlFor="true" className="flex-1 cursor-pointer text-base">
                True
              </Label>
            </div>
            <div
              className={cn(
                "flex items-center space-x-3 p-4 rounded-lg border-2 transition-all",
                selectedAnswer === "false"
                  ? "border-primary bg-primary/5"
                  : "border-border hover:border-primary/50"
              )}
            >
              <RadioGroupItem value="false" id="false" />
              <Label
                htmlFor="false"
                className="flex-1 cursor-pointer text-base"
              >
                False
              </Label>
            </div>
          </RadioGroup>
        </CardContent>
      </Card>
    );
  }

  // Render short answer question
  if (isShortAnswer) {
    return (
      <Card>
        <CardHeader>
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 bg-primary text-primary-foreground rounded-full flex items-center justify-center font-bold">
              {questionNumber}
            </div>
            <div className="flex-1">
              <h3 className="text-lg font-semibold leading-relaxed">
                {question.question_text}
              </h3>
              <p className="text-sm text-muted-foreground mt-1">
                {question.points} {question.points === 1 ? "point" : "points"}
              </p>
            </div>
          </div>
        </CardHeader>
        <CardContent>
          <Textarea
            value={selectedAnswer || ""}
            onChange={(e) => onAnswerChange(e.target.value)}
            disabled={disabled}
            placeholder="Type your answer here..."
            className="min-h-[120px]"
          />
        </CardContent>
      </Card>
    );
  }

  return null;
}
