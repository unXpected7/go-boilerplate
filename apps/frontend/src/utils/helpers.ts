import { type Schedule, type Location } from '../api/standalone';

// Utility functions for the EVV Logger

export function filterSchedulesByStatus(schedules: Schedule[], status: string): Schedule[] {
  if (status === 'all') return schedules;
  return schedules.filter(schedule => schedule.status === status);
}

export function calculateVisitDuration(startTime: Date, endTime: Date): number {
  if (!startTime || !endTime || endTime < startTime) return 0;
  const diffMs = endTime.getTime() - startTime.getTime();
  return Math.round(diffMs / (1000 * 60)); // Convert to minutes
}

export function formatTimestamp(timestamp: Date | null): string {
  if (!timestamp) return 'Invalid date';
  return timestamp.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
  });
}

export function getStatusColor(status: string): string {
  const colorMap: Record<string, string> = {
    upcoming: 'bg-blue-100 text-blue-800',
    in_progress: 'bg-yellow-100 text-yellow-800',
    completed: 'bg-green-100 text-green-800',
    missed: 'bg-red-100 text-red-800',
  };
  return colorMap[status] || 'bg-gray-100 text-gray-800';
}

export function validateGeolocation(location: Location): { isValid: boolean; error?: string } {
  if (!location || typeof location.latitude !== 'number' || typeof location.longitude !== 'number') {
    return { isValid: false, error: 'Invalid location data' };
  }

  if (location.latitude < -90 || location.latitude > 90) {
    return { isValid: false, error: 'Latitude must be between -90 and 90' };
  }

  if (location.longitude < -180 || location.longitude > 180) {
    return { isValid: false, error: 'Longitude must be between -180 and 180' };
  }

  return { isValid: true };
}

export function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: NodeJS.Timeout;
  return (...args: Parameters<T>) => {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func(...args), delay);
  };
}

export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean;
  return (...args: Parameters<T>) => {
    if (!inThrottle) {
      func(...args);
      inThrottle = true;
      setTimeout(() => (inThrottle = false), limit);
    }
  };
}