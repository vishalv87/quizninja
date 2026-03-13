# QuizNinja UI - Setup Guide

This guide will help you complete the Phase 1 setup and get the development server running.

## Current Status ✅

The following have been completed:
- ✅ Next.js project structure created
- ✅ TypeScript and Tailwind CSS configured
- ✅ All configuration files created (tsconfig.json, next.config.js, tailwind.config.ts, etc.)
- ✅ Complete folder structure set up
- ✅ Environment variables configured (.env.local and .env.example)
- ✅ ESLint and Prettier configured
- ✅ TypeScript types for all API models created
- ✅ API client with Axios and auth interceptors
- ✅ Supabase client configured
- ✅ All utility functions and constants
- ✅ Providers (Theme, Query, Toast) created
- ✅ Zustand stores (auth, UI) created
- ✅ Base layout with providers

## Next Steps 🚀

### 1. Fix NPM Permission Issue

Run this command to fix the npm cache permission issue:

```bash
sudo chown -R 501:20 "/Users/vishalvaibhav/.npm"
```

### 2. Install Dependencies

Navigate to the project directory and install all dependencies:

```bash
cd quizninja-ui
npm install
```

This will install:
- Next.js 14 with React 18
- Supabase client and auth helpers
- Axios for API calls
- TanStack Query for data fetching
- Zustand for state management
- React Hook Form + Zod for forms
- All Radix UI components (for shadcn/ui)
- Tailwind CSS and plugins
- Development tools (ESLint, Prettier, Husky)

### 3. Configure Environment Variables

Update the `.env.local` file with your actual Supabase credentials:

```bash
# Edit .env.local
nano .env.local
```

Required values:
- `NEXT_PUBLIC_SUPABASE_URL` - Your Supabase project URL
- `NEXT_PUBLIC_SUPABASE_ANON_KEY` - Your Supabase anon key

You can find these in your Supabase project settings > API.

### 4. Install shadcn/ui Components

The components.json is already configured. Install the required components:

```bash
# Install shadcn/ui components
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add input
npx shadcn-ui@latest add label
npx shadcn-ui@latest add select
npx shadcn-ui@latest add dialog
npx shadcn-ui@latest add dropdown-menu
npx shadcn-ui@latest add toast
npx shadcn-ui@latest add avatar
npx shadcn-ui@latest add badge
npx shadcn-ui@latest add tabs
npx shadcn-ui@latest add table
npx shadcn-ui@latest add skeleton
npx shadcn-ui@latest add progress
npx shadcn-ui@latest add separator
npx shadcn-ui@latest add alert
npx shadcn-ui@latest add popover
npx shadcn-ui@latest add checkbox
npx shadcn-ui@latest add radio-group
npx shadcn-ui@latest add switch
npx shadcn-ui@latest add textarea
```

Or install all at once:

```bash
npx shadcn-ui@latest add button card input label select dialog dropdown-menu toast avatar badge tabs table skeleton progress separator alert popover checkbox radio-group switch textarea
```

### 5. Start the Development Server

```bash
npm run dev
```

The application should now be running at http://localhost:3000

### 6. Verify Backend Connection

Make sure your backend API is running on port 8080:

```bash
# In the quizninja-api directory
cd ../quizninja-api
go run cmd/api/main.go
```

The frontend will connect to `http://127.0.0.1:8080/api/v1` as configured.

## Project Structure

```
quizninja-ui/
├── src/
│   ├── app/                     # Next.js App Router
│   │   ├── layout.tsx           # Root layout with providers
│   │   ├── page.tsx             # Landing page
│   │   └── globals.css          # Global styles
│   ├── components/
│   │   ├── providers/           # React providers
│   │   │   ├── theme-provider.tsx
│   │   │   ├── query-provider.tsx
│   │   │   └── toast-provider.tsx
│   │   ├── ui/                  # shadcn/ui components (to be added)
│   │   ├── auth/                # Auth components (Phase 2)
│   │   ├── quiz/                # Quiz components (Phase 3)
│   │   └── ...
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts        # Axios client with interceptors
│   │   │   └── endpoints.ts     # API endpoint constants
│   │   ├── supabase/
│   │   │   └── client.ts        # Supabase client
│   │   ├── utils.ts             # Utility functions
│   │   └── constants.ts         # App constants
│   ├── hooks/
│   │   └── useAuth.ts           # Auth hook
│   ├── store/
│   │   ├── authStore.ts         # Auth state (Zustand)
│   │   └── uiStore.ts           # UI state (Zustand)
│   ├── types/                   # TypeScript types
│   │   ├── auth.ts
│   │   ├── quiz.ts
│   │   ├── user.ts
│   │   ├── challenge.ts
│   │   ├── achievement.ts
│   │   ├── notification.ts
│   │   ├── discussion.ts
│   │   └── api.ts
│   ├── schemas/                 # Zod validation schemas (Phase 2)
│   └── config/
│       ├── env.ts               # Environment variable validation
│       └── site.ts              # Site configuration
├── public/                      # Static assets
├── .env.local                   # Environment variables (not committed)
├── .env.example                 # Environment variables template
├── components.json              # shadcn/ui configuration
├── package.json                 # Dependencies
├── tsconfig.json                # TypeScript configuration
├── tailwind.config.ts           # Tailwind configuration
├── next.config.js               # Next.js configuration
└── README.md                    # Project documentation
```

## Available Scripts

```bash
npm run dev          # Start development server
npm run build        # Build for production
npm run start        # Start production server
npm run lint         # Run ESLint
npm run type-check   # Run TypeScript compiler check
npm run format       # Format code with Prettier
npm run format:check # Check code formatting
```

## Key Features Implemented

### 1. API Client (`src/lib/api/client.ts`)
- Axios instance with base URL configuration
- Request interceptor to add Supabase auth tokens
- Response interceptor for error handling
- Automatic redirect to login on 401 errors

### 2. Supabase Client (`src/lib/supabase/client.ts`)
- Supabase initialization with auth configuration
- Helper functions for auth operations (signIn, signUp, signOut)
- Session management
- Auth state change listener

### 3. State Management
- **Auth Store** (`src/store/authStore.ts`): User authentication state with Zustand
- **UI Store** (`src/store/uiStore.ts`): UI state (sidebar, notifications)
- **TanStack Query**: Server state management (to be used in Phase 2+)

### 4. Providers
- **ThemeProvider**: Dark/light mode support with next-themes
- **QueryProvider**: TanStack Query client with DevTools
- **ToastProvider**: Toast notifications with Sonner

### 5. TypeScript Types
Complete type definitions for:
- Authentication (User, Session, Profile)
- Quizzes (Quiz, Question, QuizAttempt)
- Users (UserPreferences, UserStats, Friends)
- Challenges, Achievements, Notifications, Discussions
- API responses and errors

### 6. Utilities
- `cn()`: Class name utility (clsx + tailwind-merge)
- Date formatting functions
- Duration formatting
- String utilities
- Color helpers for difficulty and scores

## Troubleshooting

### npm install fails
Run: `sudo chown -R 501:20 "/Users/vishalvaibhav/.npm"`

### TypeScript errors
Run: `npm run type-check` to see detailed errors

### Port 3000 already in use
Kill the process: `lsof -ti:3000 | xargs kill -9`

### Backend connection fails
- Check if backend is running on port 8080
- Verify NEXT_PUBLIC_API_BASE_URL in .env.local
- Check backend ALLOWED_ORIGINS includes http://localhost:3000

## Next Phase: Authentication (Phase 2)

Once Phase 1 is complete and the dev server is running, we'll implement:
- Login page
- Registration page
- Auth middleware
- Protected routes
- Profile management

## Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [shadcn/ui Documentation](https://ui.shadcn.com/)
- [TanStack Query](https://tanstack.com/query/latest)
- [Zustand](https://github.com/pmndrs/zustand)
- [Supabase](https://supabase.com/docs)
- [Tailwind CSS](https://tailwindcss.com/docs)

---

**Phase 1 Status:** 95% Complete (pending npm install and shadcn/ui component installation)

Once you've run the npm install and added shadcn/ui components, Phase 1 will be 100% complete! 🎉