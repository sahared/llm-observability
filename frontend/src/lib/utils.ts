import { type ClassValue, clsx } from "clsx"
import { twMerge } from "tailwind-merge"

// Utility for merging Tailwind CSS classes
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Format currency (USD)
export function formatCurrency(amount: number): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 2,
    maximumFractionDigits: 6,
  }).format(amount)
}

// Format large numbers with commas
export function formatNumber(num: number): string {
  return new Intl.NumberFormat('en-US').format(num)
}

// Format duration in milliseconds to human-readable
export function formatDuration(ms: number): string {
  if (ms < 1000) {
    return `${ms}ms`
  }
  
  const seconds = ms / 1000
  if (seconds < 60) {
    return `${seconds.toFixed(2)}s`
  }
  
  const minutes = seconds / 60
  if (minutes < 60) {
    return `${minutes.toFixed(2)}m`
  }
  
  const hours = minutes / 60
  return `${hours.toFixed(2)}h`
}

// Format date to human-readable format
export function formatDate(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  
  return new Intl.DateTimeFormat('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(d)
}

// Format relative time (e.g., "2 hours ago")
export function formatRelativeTime(date: string | Date): string {
  const d = typeof date === 'string' ? new Date(date) : date
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - d.getTime()) / 1000)
  
  if (diffInSeconds < 60) {
    return `${diffInSeconds}s ago`
  }
  
  const diffInMinutes = Math.floor(diffInSeconds / 60)
  if (diffInMinutes < 60) {
    return `${diffInMinutes}m ago`
  }
  
  const diffInHours = Math.floor(diffInMinutes / 60)
  if (diffInHours < 24) {
    return `${diffInHours}h ago`
  }
  
  const diffInDays = Math.floor(diffInHours / 24)
  if (diffInDays < 30) {
    return `${diffInDays}d ago`
  }
  
  const diffInMonths = Math.floor(diffInDays / 30)
  if (diffInMonths < 12) {
    return `${diffInMonths}mo ago`
  }
  
  const diffInYears = Math.floor(diffInMonths / 12)
  return `${diffInYears}y ago`
}

// Get status color for badges
export function getStatusColor(status: string): string {
  const statusColors: Record<string, string> = {
    success: 'bg-green-100 text-green-800 border-green-200',
    error: 'bg-red-100 text-red-800 border-red-200',
    pending: 'bg-yellow-100 text-yellow-800 border-yellow-200',
    timeout: 'bg-orange-100 text-orange-800 border-orange-200',
    unknown: 'bg-gray-100 text-gray-800 border-gray-200',
  }
  
  return statusColors[status.toLowerCase()] || statusColors.unknown
}

// Get provider color
export function getProviderColor(provider: string): string {
  const providerColors: Record<string, string> = {
    openai: 'bg-emerald-100 text-emerald-800',
    anthropic: 'bg-orange-100 text-orange-800',
    cohere: 'bg-purple-100 text-purple-800',
    google: 'bg-blue-100 text-blue-800',
  }
  
  return providerColors[provider.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

// Truncate text
export function truncate(text: string, maxLength: number): string {
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

// Copy to clipboard
export async function copyToClipboard(text: string): Promise<boolean> {
  try {
    await navigator.clipboard.writeText(text)
    return true
  } catch (err) {
    console.error('Failed to copy to clipboard:', err)
    return false
  }
}
