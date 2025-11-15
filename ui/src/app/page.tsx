import Link from 'next/link'
import { Button } from '@/components/ui/button'

export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="z-10 max-w-5xl w-full items-center justify-center flex flex-col gap-8">
        <div className="text-center space-y-4">
          <h1 className="text-5xl md:text-6xl font-bold bg-gradient-to-r from-primary to-purple-600 bg-clip-text text-transparent">
            Welcome to QuizNinja
          </h1>
          <p className="text-xl md:text-2xl text-muted-foreground max-w-2xl mx-auto">
            Your journey to knowledge starts here
          </p>
          <p className="text-base text-muted-foreground max-w-xl mx-auto">
            Test your knowledge, challenge your friends, and climb the leaderboard with our engaging quiz platform.
          </p>
        </div>

        <div className="flex flex-col sm:flex-row gap-4 mt-4">
          <Button size="lg" asChild>
            <Link href="/register">
              Get Started
            </Link>
          </Button>
          <Button size="lg" variant="outline" asChild>
            <Link href="/login">
              Sign In
            </Link>
          </Button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mt-12 max-w-4xl">
          <div className="p-6 rounded-lg border bg-card text-card-foreground">
            <h3 className="font-semibold text-lg mb-2">📚 Diverse Topics</h3>
            <p className="text-sm text-muted-foreground">
              Explore quizzes across various categories and difficulty levels
            </p>
          </div>
          <div className="p-6 rounded-lg border bg-card text-card-foreground">
            <h3 className="font-semibold text-lg mb-2">🏆 Compete & Win</h3>
            <p className="text-sm text-muted-foreground">
              Challenge friends and climb the global leaderboard
            </p>
          </div>
          <div className="p-6 rounded-lg border bg-card text-card-foreground">
            <h3 className="font-semibold text-lg mb-2">🎖️ Achievements</h3>
            <p className="text-sm text-muted-foreground">
              Unlock badges and track your learning progress
            </p>
          </div>
        </div>
      </div>
    </main>
  )
}