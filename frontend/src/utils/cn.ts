/**
 * Utility for conditionally joining CSS class names
 * Filters out falsy values and joins the rest with spaces
 * 
 * @example
 * cn('base-class', isActive && 'active', 'another-class')
 * // => "base-class active another-class"
 * 
 * @param classes - Array of class names or conditional expressions
 * @returns Joined class string
 */
export function cn(...classes: (string | boolean | undefined | null)[]): string {
  return classes.filter(Boolean).join(' ');
}
