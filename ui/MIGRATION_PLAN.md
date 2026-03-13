# QuizNinja: Flutter to Next.js Migration Plan

**Version:** 1.0
**Date:** 2025-11-15
**Status:** Planning Phase

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Technology Stack](#technology-stack)
3. [Project Structure](#project-structure)
4. [Migration Strategy](#migration-strategy)
5. [Implementation Phases](#implementation-phases)
6. [Authentication & Authorization](#authentication--authorization)
7. [State Management](#state-management)
8. [API Integration](#api-integration)
9. [UI Component Library](#ui-component-library)
10. [Feature Parity Checklist](#feature-parity-checklist)
11. [Development Workflow](#development-workflow)
12. [Deployment Strategy](#deployment-strategy)
13. [Testing Strategy](#testing-strategy)
14. [Performance Considerations](#performance-considerations)

---

## Executive Summary

This document outlines the comprehensive plan for migrating the QuizNinja mobile application from Flutter to a modern Next.js web application while maintaining the existing Go backend API.

### Goals
- Build a responsive, performant web application
- Maintain 100% feature parity with Flutter app
- Improve developer experience with TypeScript
- Leverage modern React patterns and Next.js App Router
- Ensure seamless integration with existing Go backend
- Provide excellent UX with shadcn/ui components

### Key Decisions
- **Framework:** Next.js 14+ with App Router
- **Language:** TypeScript (strict mode)
- **Styling:** Tailwind CSS
- **UI Components:** shadcn/ui
- **Backend:** Existing Go API (no changes)
- **Authentication:** Supabase Auth
- **Database:** Existing PostgreSQL via Supabase
- **Deployment:** Vercel (recommended) or self-hosted

---

## Technology Stack

### Frontend Stack

| Technology | Version | Purpose |
|------------|---------|---------|
| Next.js | 14+ | React framework with App Router |
| React | 18+ | UI library |
| TypeScript | 5+ | Type-safe JavaScript |
| Tailwind CSS | 3+ | Utility-first CSS framework |
| shadcn/ui | Latest | UI component library |
| Supabase JS | Latest | Authentication client |
| Axios/Fetch | Latest | HTTP client for API calls |
| React Hook Form | Latest | Form handling |
| Zod | Latest | Schema validation |
| TanStack Query | Latest | Server state management |
| Zustand | Latest | Client state management |
| date-fns | Latest | Date utilities |
| Sonner | Latest | Toast notifications |

### Development Tools

| Tool | Purpose |
|------|---------|
| ESLint | Code linting |
| Prettier | Code formatting |
| Husky | Git hooks |
| lint-staged | Pre-commit linting |
| TypeScript | Type checking |
| pnpm/npm | Package management |

### Backend (Unchanged)

| Technology | Purpose |
|------------|---------|
| Go | Backend language |
| Gin | HTTP framework |
| PostgreSQL | Database |
| Supabase | Auth & Database hosting |

---

## Project Structure

```
quizninja-ui/
├── src/
│   ├── app/                          # Next.js App Router
│   │   ├── (auth)/                   # Auth group routes
│   │   │   ├── login/
│   │   │   │   └── page.tsx
│   │   │   ├── register/
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   ├── (dashboard)/              # Protected routes group
│   │   │   ├── dashboard/
│   │   │   │   └── page.tsx
│   │   │   ├── quizzes/
│   │   │   │   ├── page.tsx
│   │   │   │   ├── [id]/
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── take/
│   │   │   │   │       └── page.tsx
│   │   │   │   └── category/
│   │   │   │       └── [categoryId]/
│   │   │   │           └── page.tsx
│   │   │   ├── profile/
│   │   │   │   ├── page.tsx
│   │   │   │   └── edit/
│   │   │   │       └── page.tsx
│   │   │   ├── friends/
│   │   │   │   ├── page.tsx
│   │   │   │   └── requests/
│   │   │   │       └── page.tsx
│   │   │   ├── challenges/
│   │   │   │   ├── page.tsx
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx
│   │   │   ├── achievements/
│   │   │   │   └── page.tsx
│   │   │   ├── leaderboard/
│   │   │   │   └── page.tsx
│   │   │   ├── discussions/
│   │   │   │   ├── page.tsx
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx
│   │   │   ├── notifications/
│   │   │   │   └── page.tsx
│   │   │   ├── settings/
│   │   │   │   └── page.tsx
│   │   │   └── layout.tsx
│   │   ├── (onboarding)/             # Onboarding flow
│   │   │   ├── welcome/
│   │   │   │   └── page.tsx
│   │   │   └── preferences/
│   │   │       └── page.tsx
│   │   ├── api/                      # API routes (if needed)
│   │   │   └── auth/
│   │   │       └── [...nextauth]/
│   │   │           └── route.ts
│   │   ├── layout.tsx                # Root layout
│   │   ├── page.tsx                  # Landing page
│   │   └── globals.css
│   ├── components/                   # React components
│   │   ├── ui/                       # shadcn/ui components
│   │   │   ├── button.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── input.tsx
│   │   │   ├── select.tsx
│   │   │   ├── toast.tsx
│   │   │   └── ...
│   │   ├── auth/
│   │   │   ├── LoginForm.tsx
│   │   │   ├── RegisterForm.tsx
│   │   │   └── AuthGuard.tsx
│   │   ├── quiz/
│   │   │   ├── QuizCard.tsx
│   │   │   ├── QuizList.tsx
│   │   │   ├── QuestionCard.tsx
│   │   │   ├── QuizTimer.tsx
│   │   │   ├── QuizProgress.tsx
│   │   │   └── QuizResults.tsx
│   │   ├── social/
│   │   │   ├── FriendCard.tsx
│   │   │   ├── FriendList.tsx
│   │   │   ├── FriendRequestCard.tsx
│   │   │   └── UserSearch.tsx
│   │   ├── challenge/
│   │   │   ├── ChallengeCard.tsx
│   │   │   ├── ChallengeList.tsx
│   │   │   ├── CreateChallengeDialog.tsx
│   │   │   └── ChallengeResults.tsx
│   │   ├── achievement/
│   │   │   ├── AchievementCard.tsx
│   │   │   ├── AchievementGrid.tsx
│   │   │   └── AchievementToast.tsx
│   │   ├── leaderboard/
│   │   │   ├── LeaderboardTable.tsx
│   │   │   └── UserRankCard.tsx
│   │   ├── discussion/
│   │   │   ├── DiscussionCard.tsx
│   │   │   ├── DiscussionList.tsx
│   │   │   ├── ReplyCard.tsx
│   │   │   └── CreateDiscussionDialog.tsx
│   │   ├── notification/
│   │   │   ├── NotificationBell.tsx
│   │   │   ├── NotificationList.tsx
│   │   │   └── NotificationCard.tsx
│   │   ├── layout/
│   │   │   ├── Header.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Footer.tsx
│   │   │   └── MobileNav.tsx
│   │   └── common/
│   │       ├── LoadingSpinner.tsx
│   │       ├── ErrorBoundary.tsx
│   │       ├── EmptyState.tsx
│   │       └── ConfirmDialog.tsx
│   ├── lib/                          # Utility libraries
│   │   ├── api/
│   │   │   ├── client.ts             # Axios/Fetch client
│   │   │   ├── endpoints.ts          # API endpoint definitions
│   │   │   ├── auth.ts               # Auth API calls
│   │   │   ├── quiz.ts               # Quiz API calls
│   │   │   ├── user.ts               # User API calls
│   │   │   ├── friends.ts            # Friends API calls
│   │   │   ├── challenges.ts         # Challenges API calls
│   │   │   ├── achievements.ts       # Achievements API calls
│   │   │   ├── notifications.ts      # Notifications API calls
│   │   │   └── discussions.ts        # Discussions API calls
│   │   ├── supabase/
│   │   │   ├── client.ts             # Supabase client
│   │   │   └── auth.ts               # Auth helpers
│   │   ├── utils.ts                  # Utility functions
│   │   ├── cn.ts                     # Class name utility
│   │   └── constants.ts              # App constants
│   ├── hooks/                        # Custom React hooks
│   │   ├── useAuth.ts
│   │   ├── useQuiz.ts
│   │   ├── useQuizAttempt.ts
│   │   ├── useFriends.ts
│   │   ├── useChallenges.ts
│   │   ├── useAchievements.ts
│   │   ├── useNotifications.ts
│   │   ├── useLeaderboard.ts
│   │   ├── useDiscussions.ts
│   │   ├── useLocalStorage.ts
│   │   └── useDebounce.ts
│   ├── store/                        # Zustand stores
│   │   ├── authStore.ts
│   │   ├── quizStore.ts
│   │   ├── notificationStore.ts
│   │   └── uiStore.ts
│   ├── types/                        # TypeScript types
│   │   ├── api.ts                    # API response types
│   │   ├── models.ts                 # Data model types
│   │   ├── auth.ts
│   │   ├── quiz.ts
│   │   ├── user.ts
│   │   ├── challenge.ts
│   │   ├── achievement.ts
│   │   ├── notification.ts
│   │   └── discussion.ts
│   ├── schemas/                      # Zod validation schemas
│   │   ├── auth.ts
│   │   ├── quiz.ts
│   │   ├── user.ts
│   │   └── challenge.ts
│   └── config/                       # App configuration
│       ├── site.ts                   # Site metadata
│       └── env.ts                    # Environment variables
├── public/                           # Static assets
│   ├── images/
│   ├── icons/
│   └── favicon.ico
├── .env.local                        # Environment variables
├── .env.example                      # Example env file
├── .eslintrc.json                    # ESLint config
├── .prettierrc                       # Prettier config
├── next.config.js                    # Next.js config
├── tailwind.config.ts                # Tailwind config
├── tsconfig.json                     # TypeScript config
├── components.json                   # shadcn/ui config
├── package.json
├── pnpm-lock.yaml / package-lock.json
├── README.md
└── MIGRATION_PLAN.md                 # This file
```

---

## Migration Strategy

### Approach: Parallel Development

We'll build the Next.js app in parallel with the Flutter app rather than attempting a direct port. This allows us to:

1. **Optimize for Web:** Take advantage of web-specific patterns and UX
2. **Incremental Development:** Build and test features independently
3. **Risk Mitigation:** Maintain the Flutter app until web app is ready
4. **Team Efficiency:** Different team members can work on different parts

### Key Principles

1. **API-First:** Existing Go backend remains unchanged
2. **Type Safety:** Leverage TypeScript for robust code
3. **Component Reusability:** Build modular, reusable components
4. **Progressive Enhancement:** Start with core features, add enhancements
5. **Mobile-First Responsive:** Design for mobile, scale up to desktop
6. **Performance:** Optimize for Web Vitals (LCP, FID, CLS)
7. **Accessibility:** WCAG 2.1 AA compliance

---

## Implementation Phases

### Phase 1: Project Setup & Foundation (Week 1)

**Goal:** Set up the project infrastructure and development environment

#### Tasks:
1. **Initialize Next.js Project**
   ```bash
   npx create-next-app@latest quizninja-ui --typescript --tailwind --app
   cd quizninja-ui
   ```

2. **Install Core Dependencies**
   ```bash
   pnpm add @supabase/supabase-js
   pnpm add axios
   pnpm add @tanstack/react-query
   pnpm add zustand
   pnpm add react-hook-form @hookform/resolvers
   pnpm add zod
   pnpm add date-fns
   pnpm add sonner
   ```

3. **Install shadcn/ui**
   ```bash
   npx shadcn-ui@latest init
   npx shadcn-ui@latest add button card input label select dialog toast
   ```

4. **Configure Environment Variables**
   - Create `.env.local` with Supabase keys and API URL
   - Add environment variable validation

5. **Set up Project Structure**
   - Create folder structure as outlined above
   - Set up path aliases in `tsconfig.json`

6. **Configure Tools**
   - ESLint with Next.js rules
   - Prettier for code formatting
   - Husky for pre-commit hooks

7. **Create Base Layout**
   - Root layout with providers
   - Theme provider (light/dark mode)
   - Toast provider
   - TanStack Query provider

8. **Set up API Client**
   - Configure Axios instance with interceptors
   - Add authentication token handling
   - Add error handling middleware

#### Deliverables:
- ✅ Working Next.js app with all dependencies
- ✅ Environment configuration
- ✅ Base layout and providers
- ✅ API client setup

---

### Phase 2: Authentication & Authorization (Week 2)

**Goal:** Implement complete authentication flow with Supabase

#### Tasks:
1. **Supabase Client Setup**
   - Initialize Supabase client
   - Configure auth helpers for Next.js

2. **Auth Pages**
   - Login page (`/login`)
   - Registration page (`/register`)
   - Logout functionality

3. **Auth Components**
   - `LoginForm` with validation
   - `RegisterForm` with validation
   - `AuthGuard` for protected routes

4. **Auth State Management**
   - Create `authStore` with Zustand
   - Implement `useAuth` hook
   - Handle token storage and refresh

5. **Middleware Setup**
   - Create Next.js middleware for route protection
   - Redirect logic for authenticated/unauthenticated users

6. **API Integration**
   - Integrate with backend `/api/v1/auth/login`
   - Integrate with backend `/api/v1/auth/register`
   - Integrate with backend `/api/v1/auth/logout`
   - Profile fetching from `/api/v1/profile`

#### Deliverables:
- ✅ Complete auth flow (login, register, logout)
- ✅ Protected route middleware
- ✅ Auth state management
- ✅ Session persistence

---

### Phase 3: Core Quiz Features (Weeks 3-4)

**Goal:** Implement the main quiz browsing and taking functionality

#### 3.1 Quiz Browsing

**Tasks:**
1. **Quiz List Page** (`/quizzes`)
   - Display all quizzes with pagination
   - Filter by category
   - Filter by difficulty
   - Search functionality

2. **Quiz Components**
   - `QuizCard` - Display quiz summary
   - `QuizList` - Grid/list view
   - `QuizFilters` - Category/difficulty filters
   - `QuizSearch` - Search bar

3. **Quiz Detail Page** (`/quizzes/[id]`)
   - Display quiz information
   - Show questions count, time limit, points
   - Show statistics (attempts, average score)
   - Start quiz button
   - Favorite button

4. **API Integration**
   - `GET /api/v1/quizzes` - List quizzes
   - `GET /api/v1/quizzes/:id` - Quiz details
   - `GET /api/v1/quizzes/featured` - Featured quizzes
   - `GET /api/v1/quizzes/category/:id` - Category quizzes
   - `GET /api/v1/categories` - Categories list

#### 3.2 Taking Quizzes

**Tasks:**
1. **Quiz Taking Page** (`/quizzes/[id]/take`)
   - Question display with options
   - Navigation between questions
   - Timer display
   - Progress indicator
   - Submit quiz functionality
   - Save progress functionality
   - Pause/resume functionality

2. **Quiz Components**
   - `QuestionCard` - Display question and options
   - `QuizTimer` - Countdown timer
   - `QuizProgress` - Progress bar
   - `QuizNavigation` - Question navigation
   - `QuizResults` - Results display

3. **Quiz State Management**
   - Track current question
   - Track answers
   - Track time remaining
   - Handle save/pause/resume
   - Handle submission

4. **API Integration**
   - `POST /api/v1/quizzes/:id/attempts` - Start attempt
   - `PUT /api/v1/quizzes/:id/attempts/:attemptId` - Update attempt
   - `POST /api/v1/quizzes/:id/attempts/:attemptId/submit` - Submit quiz
   - `PUT /api/v1/quizzes/:id/attempts/:attemptId/save-progress` - Save progress
   - `POST /api/v1/quizzes/:id/attempts/:attemptId/pause` - Pause quiz
   - `POST /api/v1/quizzes/:id/attempts/:attemptId/resume` - Resume quiz

#### Deliverables:
- ✅ Quiz browsing with filters and search
- ✅ Quiz detail pages
- ✅ Quiz taking functionality
- ✅ Quiz results display
- ✅ Save/pause/resume functionality

---

### Phase 4: User Profile & Dashboard (Week 5)

**Goal:** Implement user profile, statistics, and dashboard

#### Tasks:
1. **Dashboard Page** (`/dashboard`)
   - User stats overview
   - Recent attempts
   - Active sessions
   - Quick actions (start quiz, view challenges)
   - Featured quizzes
   - Friend activity (if available)

2. **Profile Page** (`/profile`)
   - Display user information
   - Show achievements
   - Show statistics
   - Show attempt history
   - Edit profile button

3. **Profile Edit Page** (`/profile/edit`)
   - Update name, email
   - Update avatar
   - Update preferences

4. **User Preferences**
   - Category preferences
   - Difficulty preferences
   - Notification settings
   - Privacy settings

5. **API Integration**
   - `GET /api/v1/profile` - Get profile
   - `PUT /api/v1/profile` - Update profile
   - `GET /api/v1/users/preferences` - Get preferences
   - `PUT /api/v1/users/preferences` - Update preferences
   - `GET /api/v1/users/stats` - Get statistics
   - `GET /api/v1/users/attempts` - Get attempt history
   - `GET /api/v1/users/active-sessions` - Get active sessions

#### Deliverables:
- ✅ User dashboard
- ✅ Profile viewing and editing
- ✅ User preferences management
- ✅ Statistics display

---

### Phase 5: Social Features - Friends (Week 6)

**Goal:** Implement friends system

#### Tasks:
1. **Friends Page** (`/friends`)
   - Display friends list
   - Friend requests tab
   - Add friend button
   - User search

2. **Friend Components**
   - `FriendCard` - Display friend info
   - `FriendList` - List of friends
   - `FriendRequestCard` - Display friend request
   - `UserSearch` - Search for users

3. **Friend Management**
   - Send friend request
   - Accept/decline friend request
   - Remove friend
   - View friend profile

4. **API Integration**
   - `POST /api/v1/friends/requests` - Send request
   - `GET /api/v1/friends/requests` - Get requests
   - `PUT /api/v1/friends/requests/:id` - Accept/decline
   - `DELETE /api/v1/friends/requests/:id` - Cancel request
   - `GET /api/v1/friends` - Get friends list
   - `DELETE /api/v1/friends/:id` - Remove friend
   - `GET /api/v1/friends/search` - Search users

#### Deliverables:
- ✅ Friends list
- ✅ Friend request management
- ✅ User search
- ✅ Friend profile viewing

---

### Phase 6: Challenges System (Week 7)

**Goal:** Implement challenge creation and management

#### Tasks:
1. **Challenges Page** (`/challenges`)
   - Display active challenges
   - Display pending challenges
   - Display completed challenges
   - Create challenge button
   - Challenge statistics

2. **Challenge Components**
   - `ChallengeCard` - Display challenge info
   - `ChallengeList` - List of challenges
   - `CreateChallengeDialog` - Create challenge form
   - `ChallengeResults` - Display challenge results

3. **Challenge Management**
   - Create challenge (select friend, select quiz)
   - Accept/decline challenge
   - Take challenge quiz
   - View challenge results
   - Challenge notifications

4. **API Integration**
   - `POST /api/v1/challenges` - Create challenge
   - `GET /api/v1/challenges` - Get all challenges
   - `GET /api/v1/challenges/stats` - Get stats
   - `GET /api/v1/challenges/pending` - Get pending
   - `GET /api/v1/challenges/active` - Get active
   - `GET /api/v1/challenges/completed` - Get completed
   - `GET /api/v1/challenges/:id` - Get challenge details
   - `PUT /api/v1/challenges/:id/accept` - Accept challenge
   - `PUT /api/v1/challenges/:id/decline` - Decline challenge
   - `POST /api/v1/challenges/:id/link-attempt` - Link attempt
   - `PUT /api/v1/challenges/:id/complete` - Complete challenge

#### Deliverables:
- ✅ Challenge creation
- ✅ Challenge management
- ✅ Challenge participation
- ✅ Challenge results

---

### Phase 7: Achievements & Gamification (Week 8)

**Goal:** Implement achievements and leaderboard

#### Tasks:
1. **Achievements Page** (`/achievements`)
   - Display all achievements
   - Filter by category
   - Show locked/unlocked status
   - Show progress for locked achievements
   - Achievement statistics

2. **Leaderboard Page** (`/leaderboard`)
   - Display global leaderboard
   - Show user rank
   - Filter by time period (if available)
   - Show achievement counts

3. **Achievement Components**
   - `AchievementCard` - Display achievement
   - `AchievementGrid` - Grid layout
   - `AchievementToast` - Achievement unlock notification
   - `LeaderboardTable` - Leaderboard display
   - `UserRankCard` - User rank display

4. **Achievement System**
   - Display achievements with progress
   - Show unlock animations
   - Toast notifications for new achievements
   - Points tracking

5. **API Integration**
   - `GET /api/v1/achievements` - Get achievements
   - `GET /api/v1/achievements/progress` - Get progress
   - `GET /api/v1/achievements/stats` - Get stats
   - `POST /api/v1/achievements/check` - Check for new achievements
   - `GET /api/v1/users/achievements` - Get user achievements
   - `GET /api/v1/leaderboard` - Get leaderboard
   - `GET /api/v1/leaderboard/rank` - Get user rank
   - `GET /api/v1/leaderboard/stats` - Get stats

#### Deliverables:
- ✅ Achievements display
- ✅ Achievement progress tracking
- ✅ Leaderboard
- ✅ User ranking

---

### Phase 8: Notifications System (Week 9)

**Goal:** Implement comprehensive notifications

#### Tasks:
1. **Notifications Page** (`/notifications`)
   - Display all notifications
   - Filter by type
   - Mark as read/unread
   - Delete notifications
   - Notification statistics

2. **Notification Components**
   - `NotificationBell` - Header bell icon with badge
   - `NotificationList` - List of notifications
   - `NotificationCard` - Individual notification
   - `NotificationDropdown` - Quick view dropdown

3. **Notification Types**
   - Friend requests
   - Challenge invites
   - Achievement unlocks
   - Quiz reminders
   - Discussion replies
   - System notifications

4. **Real-time Updates** (Optional)
   - Polling or WebSocket for new notifications
   - Update badge count
   - Toast notifications

5. **API Integration**
   - `GET /api/v1/notifications` - Get notifications
   - `GET /api/v1/notifications/stats` - Get stats
   - `GET /api/v1/notifications/:id` - Get notification
   - `PUT /api/v1/notifications/:id/read` - Mark as read
   - `PUT /api/v1/notifications/:id/unread` - Mark as unread
   - `PUT /api/v1/notifications/read-all` - Mark all as read
   - `DELETE /api/v1/notifications/:id` - Delete notification

#### Deliverables:
- ✅ Notification center
- ✅ Notification bell with badge
- ✅ Real-time notification updates
- ✅ Notification management

---

### Phase 9: Discussions & Community (Week 10)

**Goal:** Implement quiz discussions and community features

#### Tasks:
1. **Discussions Page** (`/discussions`)
   - Display all discussions
   - Filter by quiz
   - Sort by popularity/recent
   - Create discussion button

2. **Discussion Detail Page** (`/discussions/[id]`)
   - Display discussion content
   - Show replies
   - Reply to discussion
   - Like/unlike discussion
   - Edit/delete own discussion

3. **Discussion Components**
   - `DiscussionCard` - Display discussion
   - `DiscussionList` - List of discussions
   - `ReplyCard` - Display reply
   - `CreateDiscussionDialog` - Create discussion form
   - `ReplyForm` - Reply form

4. **Discussion Management**
   - Create discussion
   - Reply to discussion
   - Like discussion/reply
   - Edit discussion/reply
   - Delete discussion/reply

5. **API Integration**
   - `GET /api/v1/discussions` - Get discussions
   - `POST /api/v1/discussions` - Create discussion
   - `GET /api/v1/discussions/:id` - Get discussion
   - `PUT /api/v1/discussions/:id` - Update discussion
   - `DELETE /api/v1/discussions/:id` - Delete discussion
   - `PUT /api/v1/discussions/:id/like` - Like discussion
   - `GET /api/v1/discussions/:id/replies` - Get replies
   - `POST /api/v1/discussions/:id/replies` - Create reply
   - `PUT /api/v1/discussions/replies/:replyId` - Update reply
   - `DELETE /api/v1/discussions/replies/:replyId` - Delete reply
   - `PUT /api/v1/discussions/replies/:replyId/like` - Like reply

#### Deliverables:
- ✅ Discussion browsing
- ✅ Discussion creation
- ✅ Replies and likes
- ✅ Discussion management

---

### Phase 10: Additional Features (Week 11)

**Goal:** Implement remaining features

#### Tasks:
1. **Favorites System**
   - Add/remove favorites
   - View favorites list
   - Favorites on quiz cards

2. **Settings Page** (`/settings`)
   - Account settings
   - Notification preferences
   - Privacy settings
   - Theme settings (light/dark)

3. **Onboarding Flow**
   - Welcome screen
   - Category selection
   - Preferences setup
   - Skip option

4. **Search Functionality**
   - Global search
   - Search quizzes
   - Search users
   - Search discussions

5. **Categories Management**
   - Category browsing
   - Category-specific pages

#### Deliverables:
- ✅ Favorites functionality
- ✅ Settings management
- ✅ Onboarding flow
- ✅ Search features

---

### Phase 11: Polish & Optimization (Week 12)

**Goal:** Refine UI/UX and optimize performance

#### Tasks:
1. **UI/UX Refinements**
   - Responsive design testing (mobile, tablet, desktop)
   - Animation and transitions
   - Loading states
   - Error states
   - Empty states
   - Skeleton loaders

2. **Performance Optimization**
   - Image optimization
   - Code splitting
   - Lazy loading
   - Caching strategy
   - Bundle size optimization

3. **Accessibility**
   - Keyboard navigation
   - Screen reader support
   - ARIA labels
   - Color contrast
   - Focus management

4. **Error Handling**
   - Global error boundary
   - API error handling
   - Form validation errors
   - Network error handling
   - Retry mechanisms

5. **SEO**
   - Meta tags
   - Open Graph tags
   - Sitemap
   - Robots.txt

#### Deliverables:
- ✅ Responsive design
- ✅ Performance optimizations
- ✅ Accessibility compliance
- ✅ Error handling
- ✅ SEO setup

---

### Phase 12: Testing & Deployment (Week 13)

**Goal:** Test thoroughly and deploy to production

#### Tasks:
1. **Testing**
   - Manual testing of all features
   - Cross-browser testing
   - Mobile device testing
   - Performance testing
   - Security testing

2. **Documentation**
   - User documentation
   - Developer documentation
   - API integration guide
   - Deployment guide

3. **Deployment Setup**
   - Vercel configuration
   - Environment variables setup
   - Domain configuration
   - SSL setup
   - Analytics setup

4. **Production Deployment**
   - Deploy to staging
   - QA testing
   - Deploy to production
   - Monitor for issues

#### Deliverables:
- ✅ Comprehensive testing
- ✅ Documentation
- ✅ Production deployment
- ✅ Monitoring setup

---

## Authentication & Authorization

### Strategy

We'll use **Supabase Auth** for authentication, maintaining consistency with the Flutter app.

#### Flow:
1. User signs up/logs in via Supabase client
2. Supabase returns JWT access token and refresh token
3. Store tokens securely (httpOnly cookies or localStorage)
4. Send access token in `Authorization: Bearer <token>` header for all API requests
5. Backend validates token with Supabase
6. Refresh token automatically when expired

### Implementation Details

#### Client-Side Auth (Supabase)
```typescript
// lib/supabase/client.ts
import { createClient } from '@supabase/supabase-js'

export const supabase = createClient(
  process.env.NEXT_PUBLIC_SUPABASE_URL!,
  process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
)
```

#### Auth Store (Zustand)
```typescript
// store/authStore.ts
interface AuthState {
  user: User | null
  session: Session | null
  isLoading: boolean
  setUser: (user: User | null) => void
  setSession: (session: Session | null) => void
  logout: () => Promise<void>
}
```

#### Auth Hook
```typescript
// hooks/useAuth.ts
export function useAuth() {
  const { user, session, isLoading } = useAuthStore()
  return { user, session, isLoading, isAuthenticated: !!user }
}
```

#### Protected Routes Middleware
```typescript
// middleware.ts
import { createMiddlewareClient } from '@supabase/auth-helpers-nextjs'
import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

export async function middleware(req: NextRequest) {
  const res = NextResponse.next()
  const supabase = createMiddlewareClient({ req, res })
  const { data: { session } } = await supabase.auth.getSession()

  // Protected routes
  const protectedPaths = ['/dashboard', '/quizzes', '/profile', '/friends', '/challenges']
  const isProtected = protectedPaths.some(path => req.nextUrl.pathname.startsWith(path))

  if (isProtected && !session) {
    return NextResponse.redirect(new URL('/login', req.url))
  }

  // Auth routes (redirect to dashboard if logged in)
  const authPaths = ['/login', '/register']
  const isAuthPath = authPaths.some(path => req.nextUrl.pathname.startsWith(path))

  if (isAuthPath && session) {
    return NextResponse.redirect(new URL('/dashboard', req.url))
  }

  return res
}
```

#### API Client with Auth
```typescript
// lib/api/client.ts
import axios from 'axios'
import { supabase } from '@/lib/supabase/client'

const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL,
})

// Add auth token to requests
apiClient.interceptors.request.use(async (config) => {
  const { data: { session } } = await supabase.auth.getSession()
  if (session?.access_token) {
    config.headers.Authorization = `Bearer ${session.access_token}`
  }
  return config
})

// Handle 401 errors
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      await supabase.auth.signOut()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
```

---

## State Management

### Strategy: Hybrid Approach

We'll use a combination of state management solutions based on the data type:

#### 1. Server State - TanStack Query (React Query)
For data fetched from the API (quizzes, users, challenges, etc.)

**Why?**
- Automatic caching
- Background refetching
- Optimistic updates
- Pagination support
- Loading/error states

**Example:**
```typescript
// hooks/useQuizzes.ts
import { useQuery } from '@tanstack/react-query'
import { getQuizzes } from '@/lib/api/quiz'

export function useQuizzes() {
  return useQuery({
    queryKey: ['quizzes'],
    queryFn: getQuizzes,
    staleTime: 5 * 60 * 1000, // 5 minutes
  })
}
```

#### 2. Client State - Zustand
For UI state and client-side data (auth, theme, notifications badge, etc.)

**Why?**
- Simple and lightweight
- No boilerplate
- TypeScript support
- DevTools integration

**Example:**
```typescript
// store/uiStore.ts
import { create } from 'zustand'

interface UIState {
  theme: 'light' | 'dark'
  sidebarOpen: boolean
  toggleTheme: () => void
  toggleSidebar: () => void
}

export const useUIStore = create<UIState>((set) => ({
  theme: 'light',
  sidebarOpen: true,
  toggleTheme: () => set((state) => ({
    theme: state.theme === 'light' ? 'dark' : 'light'
  })),
  toggleSidebar: () => set((state) => ({
    sidebarOpen: !state.sidebarOpen
  })),
}))
```

#### 3. Form State - React Hook Form
For form handling and validation

**Why?**
- Performance (uncontrolled inputs)
- Easy validation with Zod
- Built-in error handling

**Example:**
```typescript
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { loginSchema } from '@/schemas/auth'

export function LoginForm() {
  const { register, handleSubmit, formState: { errors } } = useForm({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data) => {
    // Handle login
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <input {...register('email')} />
      {errors.email && <span>{errors.email.message}</span>}
    </form>
  )
}
```

#### 4. URL State - Next.js Search Params
For filters, pagination, search queries

**Why?**
- Shareable URLs
- Browser history
- SEO benefits

**Example:**
```typescript
'use client'
import { useSearchParams, useRouter } from 'next/navigation'

export function QuizFilters() {
  const searchParams = useSearchParams()
  const router = useRouter()

  const category = searchParams.get('category')
  const difficulty = searchParams.get('difficulty')

  const updateFilters = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString())
    params.set(key, value)
    router.push(`?${params.toString()}`)
  }
}
```

---

## API Integration

### API Client Structure

```typescript
// lib/api/client.ts
import axios from 'axios'
import { supabase } from '@/lib/supabase/client'

export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://127.0.0.1:8080/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Request interceptor
apiClient.interceptors.request.use(
  async (config) => {
    const { data: { session } } = await supabase.auth.getSession()
    if (session?.access_token) {
      config.headers.Authorization = `Bearer ${session.access_token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor
apiClient.interceptors.response.use(
  (response) => response.data,
  async (error) => {
    if (error.response?.status === 401) {
      await supabase.auth.signOut()
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)
```

### API Modules

#### Quiz API
```typescript
// lib/api/quiz.ts
import { apiClient } from './client'
import type { Quiz, QuizAttempt, QuizFilters } from '@/types/quiz'

export const quizApi = {
  getQuizzes: (filters?: QuizFilters) =>
    apiClient.get<Quiz[]>('/quizzes', { params: filters }),

  getQuiz: (id: string) =>
    apiClient.get<Quiz>(`/quizzes/${id}`),

  getFeaturedQuizzes: () =>
    apiClient.get<Quiz[]>('/quizzes/featured'),

  startAttempt: (quizId: string) =>
    apiClient.post<QuizAttempt>(`/quizzes/${quizId}/attempts`),

  submitAttempt: (quizId: string, attemptId: string, answers: any) =>
    apiClient.post(`/quizzes/${quizId}/attempts/${attemptId}/submit`, { answers }),

  saveProgress: (quizId: string, attemptId: string, answers: any) =>
    apiClient.put(`/quizzes/${quizId}/attempts/${attemptId}/save-progress`, { answers }),
}
```

### React Query Hooks

```typescript
// hooks/useQuiz.ts
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { quizApi } from '@/lib/api/quiz'

export function useQuizzes(filters?: QuizFilters) {
  return useQuery({
    queryKey: ['quizzes', filters],
    queryFn: () => quizApi.getQuizzes(filters),
  })
}

export function useQuiz(id: string) {
  return useQuery({
    queryKey: ['quiz', id],
    queryFn: () => quizApi.getQuiz(id),
    enabled: !!id,
  })
}

export function useStartQuizAttempt() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (quizId: string) => quizApi.startAttempt(quizId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user-attempts'] })
    },
  })
}
```

---

## UI Component Library

### shadcn/ui Setup

We'll use **shadcn/ui** for our component library. It's built on Radix UI primitives with Tailwind CSS styling.

#### Installation
```bash
npx shadcn-ui@latest init
```

#### Core Components to Install
```bash
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

### Custom Components

#### QuizCard Component
```typescript
// components/quiz/QuizCard.tsx
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import type { Quiz } from '@/types/quiz'

interface QuizCardProps {
  quiz: Quiz
  onStart: () => void
}

export function QuizCard({ quiz, onStart }: QuizCardProps) {
  return (
    <Card className="hover:shadow-lg transition-shadow">
      <CardHeader>
        <h3 className="text-xl font-bold">{quiz.title}</h3>
        <p className="text-sm text-muted-foreground">{quiz.description}</p>
      </CardHeader>
      <CardContent>
        <div className="flex gap-2 flex-wrap">
          <Badge variant="secondary">{quiz.category}</Badge>
          <Badge variant="outline">{quiz.difficulty}</Badge>
          <Badge>{quiz.question_count} questions</Badge>
        </div>
      </CardContent>
      <CardFooter>
        <Button onClick={onStart} className="w-full">Start Quiz</Button>
      </CardFooter>
    </Card>
  )
}
```

### Theming

#### Tailwind Configuration
```typescript
// tailwind.config.ts
import type { Config } from 'tailwindcss'

const config: Config = {
  darkMode: ['class'],
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        // ... more colors
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
}
export default config
```

---

## Feature Parity Checklist

### Authentication ✓
- [ ] Login
- [ ] Register
- [ ] Logout
- [ ] Profile viewing
- [ ] Profile editing
- [ ] Password reset (if implemented in Flutter)

### Quizzes ✓
- [ ] Browse quizzes
- [ ] Filter by category
- [ ] Filter by difficulty
- [ ] Search quizzes
- [ ] View quiz details
- [ ] Featured quizzes
- [ ] Take quiz
- [ ] Multiple choice questions
- [ ] True/false questions
- [ ] Short answer questions
- [ ] Quiz timer
- [ ] Save progress
- [ ] Pause/resume quiz
- [ ] Submit quiz
- [ ] View results
- [ ] View attempt history
- [ ] Retry quiz

### User Profile ✓
- [ ] View profile
- [ ] Edit profile
- [ ] View statistics
- [ ] View achievements
- [ ] View attempt history
- [ ] Active quiz sessions
- [ ] User preferences
- [ ] Onboarding flow

### Friends ✓
- [ ] View friends list
- [ ] Send friend request
- [ ] Accept friend request
- [ ] Decline friend request
- [ ] Remove friend
- [ ] Search users
- [ ] View friend profile
- [ ] Friend notifications

### Challenges ✓
- [ ] View challenges
- [ ] Create challenge
- [ ] Accept challenge
- [ ] Decline challenge
- [ ] Cancel challenge
- [ ] Take challenge quiz
- [ ] View challenge results
- [ ] Challenge statistics
- [ ] Pending challenges
- [ ] Active challenges
- [ ] Completed challenges

### Achievements ✓
- [ ] View all achievements
- [ ] View unlocked achievements
- [ ] View locked achievements
- [ ] Achievement progress
- [ ] Achievement categories
- [ ] Achievement notifications
- [ ] Achievement statistics

### Leaderboard ✓
- [ ] Global leaderboard
- [ ] User rank
- [ ] Leaderboard statistics
- [ ] Friend rankings (if applicable)

### Notifications ✓
- [ ] View notifications
- [ ] Filter by type
- [ ] Mark as read
- [ ] Mark as unread
- [ ] Delete notification
- [ ] Notification badge
- [ ] Real-time updates
- [ ] Notification statistics

### Discussions ✓
- [ ] View discussions
- [ ] Create discussion
- [ ] Reply to discussion
- [ ] Like discussion
- [ ] Like reply
- [ ] Edit discussion
- [ ] Delete discussion
- [ ] Edit reply
- [ ] Delete reply
- [ ] Discussion statistics

### Favorites ✓
- [ ] Add to favorites
- [ ] Remove from favorites
- [ ] View favorites list
- [ ] Check if favorited

### Settings ✓
- [ ] Account settings
- [ ] Notification preferences
- [ ] Privacy settings
- [ ] Theme settings (light/dark)
- [ ] Category preferences
- [ ] Difficulty preferences

### Additional Features ✓
- [ ] Categories browsing
- [ ] Global search
- [ ] Responsive design
- [ ] Accessibility
- [ ] Error handling
- [ ] Loading states
- [ ] Empty states

---

## Development Workflow

### Getting Started

```bash
# Clone the repository (or create new project)
git clone <repo-url>
cd quizninja/quizninja-ui

# Install dependencies
pnpm install

# Set up environment variables
cp .env.example .env.local
# Edit .env.local with your Supabase and API credentials

# Run development server
pnpm dev

# Open http://localhost:3000
```

### Development Scripts

```json
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "type-check": "tsc --noEmit",
    "format": "prettier --write \"src/**/*.{ts,tsx,md}\"",
    "format:check": "prettier --check \"src/**/*.{ts,tsx,md}\""
  }
}
```

### Git Workflow

1. **Create feature branch** from `main`
   ```bash
   git checkout -b feature/quiz-taking
   ```

2. **Make changes** and commit
   ```bash
   git add .
   git commit -m "feat: implement quiz taking functionality"
   ```

3. **Push to remote**
   ```bash
   git push origin feature/quiz-taking
   ```

4. **Create Pull Request** for review

5. **Merge to main** after approval

### Commit Message Convention

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

---

## Deployment Strategy

### Recommended: Vercel

Vercel is the recommended platform for Next.js deployments.

#### Setup

1. **Install Vercel CLI**
   ```bash
   pnpm add -g vercel
   ```

2. **Login to Vercel**
   ```bash
   vercel login
   ```

3. **Deploy**
   ```bash
   vercel
   ```

#### Environment Variables

Configure in Vercel dashboard:
- `NEXT_PUBLIC_SUPABASE_URL`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`
- `NEXT_PUBLIC_API_BASE_URL`

#### Automatic Deployments

1. Connect GitHub repository to Vercel
2. Configure production branch (e.g., `main`)
3. Enable automatic deployments
4. Preview deployments for PRs

### Alternative: Self-Hosted

#### Docker Deployment

```dockerfile
# Dockerfile
FROM node:20-alpine AS base

# Install dependencies only when needed
FROM base AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Build the application
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm install -g pnpm && pnpm build

# Production image
FROM base AS runner
WORKDIR /app
ENV NODE_ENV production
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs
EXPOSE 3000
ENV PORT 3000

CMD ["node", "server.js"]
```

```bash
# Build and run
docker build -t quizninja-ui .
docker run -p 3000:3000 quizninja-ui
```

### Production Checklist

- [ ] Environment variables configured
- [ ] Error tracking (Sentry)
- [ ] Analytics (Google Analytics, Plausible, etc.)
- [ ] Performance monitoring
- [ ] SEO meta tags
- [ ] Sitemap generated
- [ ] robots.txt configured
- [ ] SSL certificate
- [ ] CDN configured (if needed)
- [ ] Database backups
- [ ] Monitoring alerts

---

## Testing Strategy

### Testing Levels

#### 1. Unit Testing (Optional)
Using **Vitest** or **Jest**

```bash
pnpm add -D vitest @testing-library/react @testing-library/jest-dom
```

Example:
```typescript
// components/quiz/QuizCard.test.tsx
import { render, screen } from '@testing-library/react'
import { QuizCard } from './QuizCard'

describe('QuizCard', () => {
  it('renders quiz information', () => {
    const quiz = {
      id: '1',
      title: 'Test Quiz',
      description: 'Test Description',
      category: 'Science',
      difficulty: 'Easy',
      question_count: 10,
    }

    render(<QuizCard quiz={quiz} onStart={() => {}} />)
    expect(screen.getByText('Test Quiz')).toBeInTheDocument()
  })
})
```

#### 2. Integration Testing
Test component interactions and API integration

#### 3. E2E Testing (Optional)
Using **Playwright** or **Cypress**

```bash
pnpm add -D @playwright/test
```

Example:
```typescript
// e2e/quiz.spec.ts
import { test, expect } from '@playwright/test'

test('user can take a quiz', async ({ page }) => {
  await page.goto('/quizzes/1')
  await page.click('text=Start Quiz')
  // ... interact with quiz
  await page.click('text=Submit')
  await expect(page).toHaveURL(/.*\/results/)
})
```

#### 4. Manual Testing
- Feature testing after implementation
- Cross-browser testing (Chrome, Firefox, Safari, Edge)
- Mobile device testing (iOS, Android)
- Accessibility testing (keyboard navigation, screen readers)

---

## Performance Considerations

### Next.js Optimizations

#### 1. Image Optimization
```typescript
import Image from 'next/image'

<Image
  src="/quiz-thumbnail.jpg"
  alt="Quiz Thumbnail"
  width={300}
  height={200}
  loading="lazy"
/>
```

#### 2. Code Splitting
```typescript
import dynamic from 'next/dynamic'

const QuizTaking = dynamic(() => import('@/components/quiz/QuizTaking'), {
  loading: () => <LoadingSpinner />,
})
```

#### 3. React Server Components
Use Server Components for data fetching where possible:

```typescript
// app/quizzes/page.tsx
export default async function QuizzesPage() {
  const quizzes = await getQuizzes() // Server-side fetch

  return <QuizList quizzes={quizzes} />
}
```

#### 4. Caching Strategy
- Static pages: Pre-render at build time
- Dynamic pages with revalidation: ISR (Incremental Static Regeneration)
- Client-side caching: TanStack Query

#### 5. Bundle Size Optimization
- Tree shaking
- Import only what you need
- Analyze bundle size: `pnpm add -D @next/bundle-analyzer`

### Performance Metrics

Target Web Vitals:
- **LCP** (Largest Contentful Paint): < 2.5s
- **FID** (First Input Delay): < 100ms
- **CLS** (Cumulative Layout Shift): < 0.1

---

## Risk Mitigation

### Potential Risks

1. **API Changes**
   - *Mitigation:* Maintain API versioning, create abstraction layer

2. **Authentication Issues**
   - *Mitigation:* Thorough testing, error handling, token refresh logic

3. **Performance Issues**
   - *Mitigation:* Performance monitoring, optimization, lazy loading

4. **Browser Compatibility**
   - *Mitigation:* Cross-browser testing, polyfills if needed

5. **Security Vulnerabilities**
   - *Mitigation:* Regular dependency updates, security audits

---

## Success Criteria

### Technical Criteria
- ✅ 100% feature parity with Flutter app
- ✅ All API endpoints integrated
- ✅ Authentication working correctly
- ✅ Responsive design (mobile, tablet, desktop)
- ✅ Performance metrics met (Core Web Vitals)
- ✅ Accessibility compliance (WCAG 2.1 AA)
- ✅ Cross-browser compatibility
- ✅ Error handling and edge cases covered

### User Experience Criteria
- ✅ Intuitive navigation
- ✅ Fast page loads
- ✅ Smooth animations and transitions
- ✅ Clear feedback for user actions
- ✅ Helpful error messages
- ✅ Mobile-friendly interface

---

## Timeline Summary

| Phase | Duration | Key Deliverables |
|-------|----------|------------------|
| Phase 1: Setup | Week 1 | Project infrastructure |
| Phase 2: Auth | Week 2 | Authentication flow |
| Phase 3: Quizzes | Weeks 3-4 | Quiz browsing and taking |
| Phase 4: Profile | Week 5 | User profile and dashboard |
| Phase 5: Friends | Week 6 | Friends system |
| Phase 6: Challenges | Week 7 | Challenges functionality |
| Phase 7: Achievements | Week 8 | Achievements and leaderboard |
| Phase 8: Notifications | Week 9 | Notifications system |
| Phase 9: Discussions | Week 10 | Discussion forum |
| Phase 10: Additional | Week 11 | Remaining features |
| Phase 11: Polish | Week 12 | UI/UX refinements |
| Phase 12: Testing | Week 13 | Testing and deployment |

**Total Estimated Duration:** 13 weeks (3 months)

---

## Next Steps

1. **Review this plan** with the team
2. **Set up development environment** (Phase 1)
3. **Create initial Next.js project**
4. **Begin Phase 2: Authentication**
5. **Iterate based on feedback**

---

## Appendix

### Environment Variables Template

```env
# .env.example

# Supabase
NEXT_PUBLIC_SUPABASE_URL=https://fadtsmogjmugkkugjuey.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your_supabase_anon_key

# API
NEXT_PUBLIC_API_BASE_URL=http://127.0.0.1:8080/api/v1

# App
NEXT_PUBLIC_APP_URL=http://localhost:3000
NEXT_PUBLIC_APP_NAME=QuizNinja

# Optional
NEXT_PUBLIC_SENTRY_DSN=
NEXT_PUBLIC_GA_TRACKING_ID=
```

### Useful Resources

- [Next.js Documentation](https://nextjs.org/docs)
- [shadcn/ui Documentation](https://ui.shadcn.com/)
- [TanStack Query Documentation](https://tanstack.com/query/latest)
- [Supabase Documentation](https://supabase.com/docs)
- [Tailwind CSS Documentation](https://tailwindcss.com/docs)
- [TypeScript Documentation](https://www.typescriptlang.org/docs/)

---

**Document Version:** 1.0
**Last Updated:** 2025-11-15
**Status:** Ready for Implementation