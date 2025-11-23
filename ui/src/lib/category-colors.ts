/**
 * Utility for generating consistent, dynamic colors for categories
 * Colors are derived from a hash of the category identifier, ensuring
 * the same category always gets the same color without hardcoding.
 */

// Define gradient pairs for category colors
const GRADIENT_PALETTES = [
  { from: "from-violet-500", to: "to-purple-600", bg: "bg-violet-100", text: "text-violet-700", accent: "violet" },
  { from: "from-blue-500", to: "to-cyan-600", bg: "bg-blue-100", text: "text-blue-700", accent: "blue" },
  { from: "from-emerald-500", to: "to-teal-600", bg: "bg-emerald-100", text: "text-emerald-700", accent: "emerald" },
  { from: "from-rose-500", to: "to-pink-600", bg: "bg-rose-100", text: "text-rose-700", accent: "rose" },
  { from: "from-amber-500", to: "to-orange-600", bg: "bg-amber-100", text: "text-amber-700", accent: "amber" },
  { from: "from-indigo-500", to: "to-blue-600", bg: "bg-indigo-100", text: "text-indigo-700", accent: "indigo" },
  { from: "from-fuchsia-500", to: "to-purple-600", bg: "bg-fuchsia-100", text: "text-fuchsia-700", accent: "fuchsia" },
  { from: "from-cyan-500", to: "to-sky-600", bg: "bg-cyan-100", text: "text-cyan-700", accent: "cyan" },
  { from: "from-lime-500", to: "to-green-600", bg: "bg-lime-100", text: "text-lime-700", accent: "lime" },
  { from: "from-red-500", to: "to-rose-600", bg: "bg-red-100", text: "text-red-700", accent: "red" },
] as const;

/**
 * Simple string hash function that produces consistent results
 * Uses djb2 algorithm for good distribution
 */
function hashString(str: string): number {
  let hash = 5381;
  for (let i = 0; i < str.length; i++) {
    hash = ((hash << 5) + hash) ^ str.charCodeAt(i);
  }
  return Math.abs(hash);
}

export interface CategoryColorScheme {
  /** Tailwind gradient "from" class */
  from: string;
  /** Tailwind gradient "to" class */
  to: string;
  /** Light background class for badges/icons */
  bg: string;
  /** Text color class */
  text: string;
  /** Accent name for reference */
  accent: string;
  /** Full gradient class string */
  gradient: string;
}

/**
 * Get a consistent color scheme for a category based on its ID or name
 * The same identifier will always return the same color scheme
 *
 * @param identifier - Category ID or name to generate color from
 * @returns Color scheme object with Tailwind classes
 */
export function getCategoryColor(identifier: string): CategoryColorScheme {
  const hash = hashString(identifier);
  const index = hash % GRADIENT_PALETTES.length;
  const palette = GRADIENT_PALETTES[index];

  return {
    ...palette,
    gradient: `${palette.from} ${palette.to}`,
  };
}

/**
 * Get all available color palettes
 * Useful for displaying a legend or color picker
 */
export function getAllColorPalettes() {
  return GRADIENT_PALETTES.map((palette) => ({
    ...palette,
    gradient: `${palette.from} ${palette.to}`,
  }));
}
