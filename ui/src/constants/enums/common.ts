/**
 * Common enums used across multiple domains
 */

// Sort Order
export const SortOrder = {
  ASC: 'asc',
  DESC: 'desc',
} as const;

export type SortOrder = typeof SortOrder[keyof typeof SortOrder];
export const SORT_ORDERS = Object.values(SortOrder);

// Type guard for sort order
export function isSortOrder(value: unknown): value is SortOrder {
  return typeof value === 'string' && SORT_ORDERS.includes(value as SortOrder);
}
