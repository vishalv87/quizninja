"use client";
 
 import { useEffect } from "react";
 import { Card, CardContent } from "@/components/ui/card";
 import { Badge } from "@/components/ui/badge";
 import { Button } from "@/components/ui/button";
 import { Separator } from "@/components/ui/separator";
 import type { QuizResults as QuizResultsType } from "@/types/quiz";
 import {
   Trophy,
   Target,
   Clock,
   CheckCircle2,
   XCircle,
   ArrowRight,
   Share2,
 } from "lucide-react";
 import Link from "next/link";
 import { Progress } from "@/components/ui/progress";
 import { motion } from "framer-motion";
 import confetti from "canvas-confetti";
 
 interface QuizResultsProps {
   results: QuizResultsType;
 }
 
 export function QuizResults({ results }: QuizResultsProps) {
   const { attempt, quiz, percentage, passed } = results;
 
   // Calculate stats
   const correctAnswers = attempt.score ?? 0;
   const totalQuestions = attempt.total_points ?? 0;
   const accuracy =
     totalQuestions > 0 ? (correctAnswers / totalQuestions) * 100 : 0;
   const incorrectAnswers = totalQuestions - correctAnswers;
 
   // Format time
   const timeSpentMinutes = attempt.time_spent
     ? Math.floor(attempt.time_spent / 60)
     : 0;
   const timeSpentSeconds = attempt.time_spent ? attempt.time_spent % 60 : 0;
 
   // Trigger confetti on load if passed
   useEffect(() => {
     if (passed) {
       const duration = 3 * 1000;
       const animationEnd = Date.now() + duration;
       const defaults = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 0 };
 
       const random = (min: number, max: number) =>
         Math.random() * (max - min) + min;
 
       const interval: any = setInterval(function () {
         const timeLeft = animationEnd - Date.now();
 
         if (timeLeft <= 0) {
           return clearInterval(interval);
         }
 
         const particleCount = 50 * (timeLeft / duration);
         confetti({
           ...defaults,
           particleCount,
           origin: { x: random(0.1, 0.3), y: Math.random() - 0.2 },
         });
         confetti({
           ...defaults,
           particleCount,
           origin: { x: random(0.7, 0.9), y: Math.random() - 0.2 },
         });
       }, 250);
 
       return () => clearInterval(interval);
     }
   }, [passed]);
 
   const container = {
     hidden: { opacity: 0 },
     show: {
       opacity: 1,
       transition: {
         staggerChildren: 0.1,
       },
     },
   };
 
   const item = {
     hidden: { opacity: 0, y: 20 },
     show: { opacity: 1, y: 0 },
   };
 
   return (
     <motion.div
       variants={container}
       initial="hidden"
       animate="show"
       className="space-y-8"
     >
       {/* Hero Section */}
       <motion.div variants={item} className="relative overflow-hidden rounded-3xl">
         <div
           className={`absolute inset-0 opacity-10 ${
             passed
               ? "bg-gradient-to-br from-green-500 via-emerald-500 to-teal-500"
               : "bg-gradient-to-br from-orange-500 via-red-500 to-pink-500"
           }`}
         />
         <div
           className={`absolute inset-0 backdrop-blur-3xl ${
             passed ? "bg-green-500/5" : "bg-orange-500/5"
           }`}
         />
         
         <div className="relative p-8 md:p-12 text-center space-y-6">
           <motion.div
             initial={{ scale: 0 }}
             animate={{ scale: 1 }}
             transition={{
               type: "spring",
               stiffness: 260,
               damping: 20,
               delay: 0.2,
             }}
             className="mx-auto w-24 h-24 rounded-full bg-background shadow-xl flex items-center justify-center mb-6"
           >
             {passed ? (
               <Trophy className="w-12 h-12 text-green-500" />
             ) : (
               <Target className="w-12 h-12 text-orange-500" />
             )}
           </motion.div>
 
           <div className="space-y-2">
             <h1 className="text-4xl md:text-5xl font-bold tracking-tight">
               {passed ? "Excellent Work!" : "Keep Going!"}
             </h1>
             <p className="text-xl text-muted-foreground max-w-lg mx-auto">
               {passed
                 ? "You've mastered this quiz. Great job!"
                 : "Don't give up. Review your answers and try again."}
             </p>
           </div>
 
           <div className="flex flex-col items-center gap-2">
             <div className="text-7xl md:text-8xl font-black tracking-tighter tabular-nums">
               {percentage.toFixed(0)}
               <span className="text-4xl md:text-5xl text-muted-foreground font-medium ml-1">
                 %
               </span>
             </div>
             <Badge
               variant={passed ? "default" : "destructive"}
               className="text-lg px-6 py-1.5 rounded-full shadow-lg"
             >
               {passed ? "Passed" : "Not Passed"}
             </Badge>
           </div>
         </div>
       </motion.div>
 
       {/* Stats Grid */}
       <motion.div
         variants={container}
         className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4"
       >
         <motion.div variants={item}>
           <Card className="h-full hover:shadow-md transition-shadow">
             <CardContent className="pt-6">
               <div className="flex flex-col gap-4">
                 <div className="p-3 bg-primary/10 w-fit rounded-2xl">
                   <Trophy className="h-6 w-6 text-primary" />
                 </div>
                 <div>
                   <p className="text-sm font-medium text-muted-foreground">
                     Total Score
                   </p>
                   <p className="text-3xl font-bold">
                     {attempt.score}
                     <span className="text-base text-muted-foreground font-normal ml-1">
                       / {attempt.total_points}
                     </span>
                   </p>
                 </div>
               </div>
             </CardContent>
           </Card>
         </motion.div>
 
         <motion.div variants={item}>
           <Card className="h-full hover:shadow-md transition-shadow">
             <CardContent className="pt-6">
               <div className="flex flex-col gap-4">
                 <div className="p-3 bg-green-500/10 w-fit rounded-2xl">
                   <CheckCircle2 className="h-6 w-6 text-green-600 dark:text-green-400" />
                 </div>
                 <div>
                   <p className="text-sm font-medium text-muted-foreground">
                     Correct Answers
                   </p>
                   <p className="text-3xl font-bold">{correctAnswers}</p>
                 </div>
               </div>
             </CardContent>
           </Card>
         </motion.div>
 
         <motion.div variants={item}>
           <Card className="h-full hover:shadow-md transition-shadow">
             <CardContent className="pt-6">
               <div className="flex flex-col gap-4">
                 <div className="p-3 bg-red-500/10 w-fit rounded-2xl">
                   <XCircle className="h-6 w-6 text-red-600 dark:text-red-400" />
                 </div>
                 <div>
                   <p className="text-sm font-medium text-muted-foreground">
                     Incorrect Answers
                   </p>
                   <p className="text-3xl font-bold">{incorrectAnswers}</p>
                 </div>
               </div>
             </CardContent>
           </Card>
         </motion.div>
 
         <motion.div variants={item}>
           <Card className="h-full hover:shadow-md transition-shadow">
             <CardContent className="pt-6">
               <div className="flex flex-col gap-4">
                 <div className="p-3 bg-blue-500/10 w-fit rounded-2xl">
                   <Clock className="h-6 w-6 text-blue-600 dark:text-blue-400" />
                 </div>
                 <div>
                   <p className="text-sm font-medium text-muted-foreground">
                     Time Spent
                   </p>
                   <p className="text-3xl font-bold">
                     {timeSpentMinutes}
                     <span className="text-base font-normal text-muted-foreground">
                       m
                     </span>{" "}
                     {timeSpentSeconds}
                     <span className="text-base font-normal text-muted-foreground">
                       s
                     </span>
                   </p>
                 </div>
               </div>
             </CardContent>
           </Card>
         </motion.div>
       </motion.div>
 
       {/* Quiz Details & Actions */}
       <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
         <motion.div variants={item} className="lg:col-span-2">
           <Card className="h-full">
             <CardContent className="p-6 md:p-8 space-y-8">
               <div>
                 <h3 className="text-xl font-semibold mb-1">Performance Analysis</h3>
                 <p className="text-muted-foreground">
                   Detailed breakdown of your quiz performance
                 </p>
               </div>
 
               <div className="space-y-6">
                 <div className="space-y-2">
                   <div className="flex justify-between text-sm">
                     <span className="font-medium">Accuracy</span>
                     <span className="text-muted-foreground">
                       {accuracy.toFixed(1)}%
                     </span>
                   </div>
                   <Progress
                     value={accuracy}
                     className={`h-3 ${
                       passed
                         ? "bg-green-100 dark:bg-green-900 [&>div]:bg-green-500"
                         : "bg-orange-100 dark:bg-orange-900 [&>div]:bg-orange-500"
                     }`}
                   />
                 </div>
 
                 <div className="grid grid-cols-2 gap-4 pt-4">
                   <div className="space-y-1">
                     <span className="text-xs text-muted-foreground uppercase tracking-wider">
                       Category
                     </span>
                     <p className="font-medium">{quiz.category}</p>
                   </div>
                   <div className="space-y-1">
                     <span className="text-xs text-muted-foreground uppercase tracking-wider">
                       Difficulty
                     </span>
                     <Badge variant="secondary" className="capitalize">
                       {quiz.difficulty}
                     </Badge>
                   </div>
                 </div>
               </div>
             </CardContent>
           </Card>
         </motion.div>
 
         <motion.div variants={item} className="space-y-4">
           <Card className="bg-primary/5 border-primary/20">
             <CardContent className="p-6">
               <h3 className="font-semibold mb-4">What's Next?</h3>
               <div className="space-y-3">
                 <Link href="/quizzes" className="block">
                   <Button className="w-full gap-2" size="lg">
                     Browse More Quizzes
                     <ArrowRight className="w-4 h-4" />
                   </Button>
                 </Link>
                 <Button variant="ghost" className="w-full gap-2">
                   <Share2 className="w-4 h-4" />
                   Share Result
                 </Button>
               </div>
             </CardContent>
           </Card>
         </motion.div>
       </div>
     </motion.div>
   );
 }
