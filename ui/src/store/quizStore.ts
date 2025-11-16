import { create } from "zustand";
import type { QuizAnswer, Quiz, QuizAttempt } from "@/types/quiz";

interface QuizState {
  // Current quiz data
  currentQuiz: Quiz | null;
  currentAttempt: QuizAttempt | null;

  // Quiz taking state
  currentQuestionIndex: number;
  answers: Record<string, QuizAnswer>; // keyed by question_id
  timeRemaining: number | null; // seconds
  isPaused: boolean;
  startTime: number | null;

  // Actions
  setCurrentQuiz: (quiz: Quiz) => void;
  setCurrentAttempt: (attempt: QuizAttempt) => void;
  setCurrentQuestionIndex: (index: number) => void;
  nextQuestion: () => void;
  previousQuestion: () => void;
  goToQuestion: (index: number) => void;

  // Answer management
  setAnswer: (questionId: string, answer: QuizAnswer) => void;
  getAnswer: (questionId: string) => QuizAnswer | undefined;
  clearAnswers: () => void;

  // Timer management
  setTimeRemaining: (seconds: number) => void;
  decrementTime: () => void;
  setPaused: (paused: boolean) => void;

  // Quiz lifecycle
  startQuiz: (quiz: Quiz, attempt: QuizAttempt) => void;
  resetQuiz: () => void;
}

export const useQuizStore = create<QuizState>((set, get) => ({
  // Initial state
  currentQuiz: null,
  currentAttempt: null,
  currentQuestionIndex: 0,
  answers: {},
  timeRemaining: null,
  isPaused: false,
  startTime: null,

  // Set current quiz
  setCurrentQuiz: (quiz) => set({ currentQuiz: quiz }),

  // Set current attempt
  setCurrentAttempt: (attempt) => set({ currentAttempt: attempt }),

  // Set current question index
  setCurrentQuestionIndex: (index) => set({ currentQuestionIndex: index }),

  // Go to next question
  nextQuestion: () => {
    const { currentQuiz, currentQuestionIndex } = get();
    if (currentQuiz && currentQuestionIndex < currentQuiz.question_count - 1) {
      set({ currentQuestionIndex: currentQuestionIndex + 1 });
    }
  },

  // Go to previous question
  previousQuestion: () => {
    const { currentQuestionIndex } = get();
    if (currentQuestionIndex > 0) {
      set({ currentQuestionIndex: currentQuestionIndex - 1 });
    }
  },

  // Go to specific question
  goToQuestion: (index) => {
    const { currentQuiz } = get();
    if (currentQuiz && index >= 0 && index < currentQuiz.question_count) {
      set({ currentQuestionIndex: index });
    }
  },

  // Set answer for a question
  setAnswer: (questionId, answer) => {
    set((state) => ({
      answers: {
        ...state.answers,
        [questionId]: answer,
      },
    }));
  },

  // Get answer for a question
  getAnswer: (questionId) => {
    return get().answers[questionId];
  },

  // Clear all answers
  clearAnswers: () => set({ answers: {} }),

  // Set time remaining
  setTimeRemaining: (seconds) => set({ timeRemaining: seconds }),

  // Decrement time by 1 second
  decrementTime: () => {
    const { timeRemaining } = get();
    if (timeRemaining !== null && timeRemaining > 0) {
      set({ timeRemaining: timeRemaining - 1 });
    }
  },

  // Set paused state
  setPaused: (paused) => set({ isPaused: paused }),

  // Start quiz - initialize all state
  startQuiz: (quiz, attempt) => {
    const timeLimit = quiz.time_limit
      ? quiz.time_limit * 60
      : null;

    set({
      currentQuiz: quiz,
      currentAttempt: attempt,
      currentQuestionIndex: 0,
      answers: {},
      timeRemaining: timeLimit,
      isPaused: false,
      startTime: Date.now(),
    });
  },

  // Reset quiz - clear all state
  resetQuiz: () => {
    set({
      currentQuiz: null,
      currentAttempt: null,
      currentQuestionIndex: 0,
      answers: {},
      timeRemaining: null,
      isPaused: false,
      startTime: null,
    });
  },
}));
