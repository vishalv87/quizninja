/**
 * Discussion-related enums and constants
 * Single source of truth for discussion sort options and types
 */

// Discussion Sort Options
export const DiscussionSort = {
  RECENT: 'recent',
  POPULAR: 'popular',
} as const;

export type DiscussionSort = typeof DiscussionSort[keyof typeof DiscussionSort];
export const DISCUSSION_SORTS = Object.values(DiscussionSort);

// Type guard for discussion sort
export function isDiscussionSort(value: unknown): value is DiscussionSort {
  return typeof value === 'string' && DISCUSSION_SORTS.includes(value as DiscussionSort);
}

// Discussion Types
export const DiscussionType = {
  QUESTION: 'question',
  GENERAL: 'general',
  BUG_REPORT: 'bug_report',
  FEATURE_REQUEST: 'feature_request',
  DISCUSSION: 'discussion',
} as const;

export type DiscussionType = typeof DiscussionType[keyof typeof DiscussionType];
export const DISCUSSION_TYPES = Object.values(DiscussionType);

// Type guard for discussion type
export function isDiscussionType(value: unknown): value is DiscussionType {
  return typeof value === 'string' && DISCUSSION_TYPES.includes(value as DiscussionType);
}
