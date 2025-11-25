# Dashboard Modernization Plan - Neumorphic UI

## Overview
Modernize the QuizNinja dashboard with a **Neumorphic design** featuring:
- Soft UI with inset/outset shadows creating 3D depth
- Floating/detached sidebar
- Violet/Indigo color palette (refined)
- Micro-animations, data visualization, better hierarchy, and interactive cards

---

## Files to Modify

| File | Purpose |
|------|---------|
| `src/app/globals.css` | Neumorphic CSS variables & utilities |
| `tailwind.config.ts` | Custom shadows, animations |
| `src/app/(dashboard)/layout.tsx` | Floating layout structure |
| `src/components/layout/Header.tsx` | Floating top bar |
| `src/components/layout/Sidebar.tsx` | Detached floating sidebar |
| `src/components/layout/MobileNav.tsx` | Match neumorphic style |
| `src/app/(dashboard)/dashboard/page.tsx` | All dashboard sections |
| `src/components/dashboard/RecentActivity.tsx` | Activity feed redesign |
| `src/components/dashboard/FeaturedQuizzesDashboard.tsx` | Featured quizzes redesign |

---

## Step 1: Global Styles (`globals.css`)

Add neumorphic CSS variables and utility classes:

```css
@layer base {
  :root {
    /* Existing variables... */

    /* Neumorphic surface */
    --neu-bg: 245 15% 95%;

    /* Shadows */
    --neu-shadow-distance: 6px;
    --neu-shadow-blur: 12px;
    --neu-shadow-light: rgba(255, 255, 255, 0.8);
    --neu-shadow-dark: rgba(163, 177, 198, 0.6);

    /* Violet accent refinement */
    --violet-glow: 263 70% 50%;
  }

  .dark {
    --neu-bg: 230 25% 10%;
    --neu-shadow-light: rgba(45, 50, 70, 0.5);
    --neu-shadow-dark: rgba(5, 5, 15, 0.8);
  }
}

@layer utilities {
  .neu-flat {
    background: hsl(var(--neu-bg));
    box-shadow:
      var(--neu-shadow-distance) var(--neu-shadow-distance) var(--neu-shadow-blur) var(--neu-shadow-dark),
      calc(var(--neu-shadow-distance) * -1) calc(var(--neu-shadow-distance) * -1) var(--neu-shadow-blur) var(--neu-shadow-light);
  }

  .neu-convex {
    background: linear-gradient(145deg,
      hsl(var(--neu-bg) / 1.05),
      hsl(var(--neu-bg) / 0.95));
    box-shadow:
      var(--neu-shadow-distance) var(--neu-shadow-distance) var(--neu-shadow-blur) var(--neu-shadow-dark),
      calc(var(--neu-shadow-distance) * -1) calc(var(--neu-shadow-distance) * -1) var(--neu-shadow-blur) var(--neu-shadow-light);
  }

  .neu-concave {
    background: linear-gradient(145deg,
      hsl(var(--neu-bg) / 0.95),
      hsl(var(--neu-bg) / 1.05));
    box-shadow:
      inset var(--neu-shadow-distance) var(--neu-shadow-distance) var(--neu-shadow-blur) var(--neu-shadow-dark),
      inset calc(var(--neu-shadow-distance) * -1) calc(var(--neu-shadow-distance) * -1) var(--neu-shadow-blur) var(--neu-shadow-light);
  }

  .neu-pressed {
    box-shadow:
      inset 2px 2px 5px var(--neu-shadow-dark),
      inset -2px -2px 5px var(--neu-shadow-light);
  }
}
```

---

## Step 2: Tailwind Config (`tailwind.config.ts`)

Add custom shadows and animations:

```typescript
theme: {
  extend: {
    boxShadow: {
      'neu': '6px 6px 12px var(--neu-shadow-dark), -6px -6px 12px var(--neu-shadow-light)',
      'neu-sm': '3px 3px 6px var(--neu-shadow-dark), -3px -3px 6px var(--neu-shadow-light)',
      'neu-inset': 'inset 4px 4px 8px var(--neu-shadow-dark), inset -4px -4px 8px var(--neu-shadow-light)',
      'violet-glow': '0 0 20px rgba(139, 92, 246, 0.3)',
    },
    keyframes: {
      'float': {
        '0%, 100%': { transform: 'translateY(0)' },
        '50%': { transform: 'translateY(-5px)' },
      },
      'pulse-soft': {
        '0%, 100%': { opacity: '1' },
        '50%': { opacity: '0.8' },
      },
      'count-up': {
        '0%': { opacity: '0', transform: 'translateY(10px)' },
        '100%': { opacity: '1', transform: 'translateY(0)' },
      },
    },
    animation: {
      'float': 'float 3s ease-in-out infinite',
      'pulse-soft': 'pulse-soft 2s ease-in-out infinite',
      'count-up': 'count-up 0.5s ease-out',
    },
  },
},
```

---

## Step 3: Dashboard Layout (`layout.tsx`)

Create floating layout structure:

```tsx
<div className="min-h-screen bg-[hsl(var(--neu-bg))]">
  <Header />
  <div className="flex flex-1 p-4 gap-4">
    <Sidebar />
    <MobileNav />
    <main className="flex-1 overflow-y-auto">
      <div className="container mx-auto px-2 py-4">
        {children}
      </div>
    </main>
  </div>
</div>
```

---

## Step 4: Header (`Header.tsx`)

Floating neumorphic header:

- Add `mx-4 mt-4 rounded-2xl` for floating effect
- Apply `neu-convex` class to header
- Inset search bar with `neu-concave`
- Icon buttons with hover scale and press animations

```tsx
<header className="sticky top-4 z-40 mx-4 rounded-2xl neu-convex">
  {/* Search bar */}
  <Button className="neu-concave rounded-xl">
    <Search />
    Search...
  </Button>

  {/* Action buttons with hover effects */}
  <Button className="rounded-xl hover:scale-105 active:neu-pressed transition-all">
    <Sun />
  </Button>
</header>
```

---

## Step 5: Sidebar (`Sidebar.tsx`)

Floating detached sidebar:

```tsx
<div className="hidden md:block w-64 ml-4 my-4 rounded-3xl neu-convex">
  <ScrollArea className="h-[calc(100vh-6rem)] py-6">
    <nav className="space-y-1.5 px-4">
      {navigation.map((item) => (
        <Link
          className={cn(
            'flex items-center gap-3 px-4 py-2.5 rounded-xl transition-all duration-300',
            isActive
              ? 'neu-concave text-violet-600'
              : 'hover:neu-pressed hover:text-violet-500'
          )}
        >
          <Icon className="h-5 w-5 transition-transform hover:scale-110" />
          {item.name}
        </Link>
      ))}
    </nav>
  </ScrollArea>
</div>
```

---

## Step 6: Dashboard Page (`page.tsx`)

### Hero Section
```tsx
<div className="relative overflow-hidden rounded-3xl neu-convex p-8 lg:p-12">
  <div className="absolute inset-0 bg-gradient-to-r from-violet-500/10 to-indigo-500/10" />
  <h1 className="text-3xl font-bold">Welcome back!</h1>
  <Button className="neu-convex hover:neu-pressed active:scale-95 transition-all">
    Start a Quiz
  </Button>
</div>
```

### Stats Cards with Circular Progress
```tsx
<Card className="neu-convex rounded-2xl hover:-translate-y-1 hover:shadow-violet-glow transition-all">
  <CardContent className="p-6">
    <div className="flex items-center justify-between">
      <CircularProgress value={75} size={60}>
        <Icon className="h-6 w-6 text-violet-500" />
      </CircularProgress>
      <div>
        <AnimatedNumber value={stats.total} />
        <p className="text-muted-foreground">Quizzes</p>
      </div>
    </div>
  </CardContent>
</Card>
```

### Quick Actions
```tsx
<div className="neu-flat rounded-2xl p-6 hover:neu-pressed active:scale-98 transition-all cursor-pointer">
  <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-violet-500 to-indigo-500 flex items-center justify-center">
    <Icon className="text-white h-6 w-6" />
  </div>
  <h3 className="font-semibold mt-4">{title}</h3>
</div>
```

---

## Step 7: New Components

### CircularProgress (`src/components/ui/circular-progress.tsx`)

```tsx
interface CircularProgressProps {
  value: number;      // 0-100
  size?: number;      // px
  strokeWidth?: number;
  color?: string;
  children?: ReactNode;
}

export function CircularProgress({
  value,
  size = 60,
  strokeWidth = 4,
  color = 'stroke-violet-500',
  children
}: CircularProgressProps) {
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (value / 100) * circumference;

  return (
    <div className="relative" style={{ width: size, height: size }}>
      <svg className="transform -rotate-90" width={size} height={size}>
        <circle
          className="stroke-muted"
          strokeWidth={strokeWidth}
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
        />
        <circle
          className={cn(color, 'transition-all duration-500')}
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          fill="transparent"
          r={radius}
          cx={size / 2}
          cy={size / 2}
          style={{ strokeDasharray: circumference, strokeDashoffset: offset }}
        />
      </svg>
      <div className="absolute inset-0 flex items-center justify-center">
        {children}
      </div>
    </div>
  );
}
```

### AnimatedNumber (`src/components/ui/animated-number.tsx`)

```tsx
'use client';

import { useEffect, useState } from 'react';

interface AnimatedNumberProps {
  value: number;
  duration?: number;
  prefix?: string;
  suffix?: string;
}

export function AnimatedNumber({
  value,
  duration = 500,
  prefix = '',
  suffix = ''
}: AnimatedNumberProps) {
  const [displayValue, setDisplayValue] = useState(0);

  useEffect(() => {
    const startTime = Date.now();
    const startValue = displayValue;

    const animate = () => {
      const elapsed = Date.now() - startTime;
      const progress = Math.min(elapsed / duration, 1);
      const eased = 1 - Math.pow(1 - progress, 3); // ease-out cubic

      setDisplayValue(Math.round(startValue + (value - startValue) * eased));

      if (progress < 1) {
        requestAnimationFrame(animate);
      }
    };

    requestAnimationFrame(animate);
  }, [value, duration]);

  return (
    <span className="tabular-nums">
      {prefix}{displayValue.toLocaleString()}{suffix}
    </span>
  );
}
```

---

## Step 8: RecentActivity (`RecentActivity.tsx`)

```tsx
<Card className="neu-flat rounded-3xl border-none">
  <CardHeader>
    <CardTitle className="flex items-center gap-2">
      <Clock className="h-5 w-5 text-violet-500" />
      Recent Activity
    </CardTitle>
  </CardHeader>
  <CardContent>
    <div className="space-y-3">
      {attempts.map((attempt, index) => (
        <motion.div
          key={attempt.id}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: index * 0.05 }}
          className="neu-concave rounded-xl p-4 flex items-center justify-between"
        >
          {/* Content */}
          <div className="flex items-center gap-4">
            <CircularProgress value={percentage} size={48}>
              <span className="text-xs font-bold">{percentage}%</span>
            </CircularProgress>
            <div>
              <p className="font-medium">{attempt.quiz_title}</p>
              <p className="text-sm text-muted-foreground">{timeAgo}</p>
            </div>
          </div>
          <Badge>{attempt.status}</Badge>
        </motion.div>
      ))}
    </div>
  </CardContent>
</Card>
```

---

## Step 9: FeaturedQuizzes (`FeaturedQuizzesDashboard.tsx`)

```tsx
<Card className="neu-flat rounded-3xl border-none h-full">
  <CardHeader>
    <CardTitle className="flex items-center gap-2">
      <Star className="h-5 w-5 fill-yellow-400 text-yellow-400" />
      Featured Quizzes
    </CardTitle>
  </CardHeader>
  <CardContent>
    <div className="space-y-4">
      {quizzes.map((quiz) => (
        <Link key={quiz.id} href={`/quizzes/${quiz.id}`}>
          <div className="neu-convex rounded-xl p-4 hover:-translate-y-1 hover:shadow-violet-glow transition-all">
            <div className="flex gap-4">
              <div className="h-16 w-16 rounded-lg bg-gradient-to-br from-violet-100 to-indigo-100 flex items-center justify-center">
                <span className="text-violet-500 font-bold text-xl">
                  {quiz.title.charAt(0)}
                </span>
              </div>
              <div className="flex-1">
                <h4 className="font-semibold">{quiz.title}</h4>
                <p className="text-sm text-muted-foreground line-clamp-1">
                  {quiz.description}
                </p>
                <div className="flex gap-2 mt-2">
                  <span className="neu-pressed text-xs px-2 py-0.5 rounded-full">
                    {quiz.category}
                  </span>
                  <span className="text-xs text-muted-foreground">
                    {quiz.questions_count} questions
                  </span>
                </div>
              </div>
            </div>
          </div>
        </Link>
      ))}
    </div>
  </CardContent>
</Card>
```

---

## Color Palette

| Use Case | Light Mode | Dark Mode |
|----------|-----------|-----------|
| Background | `#EDEEF2` | `#141620` |
| Surface | `#E3E5EB` | `#1C1F2E` |
| Primary accent | `#8B5CF6` | `#A78BFA` |
| Secondary accent | `#6366F1` | `#818CF8` |
| Text primary | `#1E1B4B` | `#F5F3FF` |
| Text muted | `#6B7280` | `#9CA3AF` |

---

## Animation Details (Framer Motion)

### Micro-interactions:
1. **Button press:** `scale: 0.95` on tap
2. **Card hover:** `y: -4` with shadow increase
3. **Icon hover:** `scale: 1.1, rotate: 5deg`
4. **Page load:** Staggered children with `fadeIn + slideUp`
5. **Number count:** Animate from 0 to value over 500ms

### Page transitions:
- Stats cards: stagger delay 0.1s each
- Activity items: stagger delay 0.05s each

---

## Implementation Order

1. Global styles and Tailwind config (foundation)
2. Dashboard layout (floating structure)
3. Header (floating top bar)
4. Sidebar (floating nav)
5. New UI components (CircularProgress, AnimatedNumber)
6. Dashboard page sections (hero, stats, actions)
7. Recent Activity and Featured Quizzes refinements
8. Mobile navigation updates
9. Final polish and animation tuning
