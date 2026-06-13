import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

/**
 * Merge Tailwind CSS classes with clsx
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

/**
 * Format price to currency string
 */
export function formatPrice(price: number, currency: string = 'IDR'): string {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(price);
}

/**
 * Format date string
 */
export function formatDate(dateString: string): string {
  return new Intl.DateTimeFormat('id-ID', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateString));
}

/**
 * Truncate text with ellipsis
 */
export function truncateText(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text;
  return text.slice(0, maxLength) + '...';
}

/**
 * Get status badge color
 */
export function getStatusColor(status: string): string {
  switch (status.toLowerCase()) {
    case 'paid':
      return 'bg-emerald-500/15 text-emerald-400';
    case 'pending':
      return 'bg-amber-500/15 text-amber-400';
    case 'failed':
      return 'bg-red-500/15 text-red-400';
    default:
      return 'bg-white/[0.06] text-slate-400';
  }
}
