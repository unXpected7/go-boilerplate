import {
  filterSchedulesByStatus,
  calculateVisitDuration,
  formatTimestamp,
  getStatusColor,
  validateGeolocation
} from '../utils/helpers';

describe('Utility Functions', () => {
  const mockSchedules = [
    { id: '1', status: 'upcoming', clientName: 'John Smith' },
    { id: '2', status: 'completed', clientName: 'Mary Johnson' },
    { id: '3', status: 'missed', clientName: 'Robert Brown' },
  ];

  describe('filterSchedulesByStatus', () => {
    it('should filter schedules by status', () => {
      const filtered = filterSchedulesByStatus(mockSchedules, 'completed');
      expect(filtered).toHaveLength(1);
      expect(filtered[0].status).toBe('completed');
    });

    it('should return all schedules when status is "all"', () => {
      const filtered = filterSchedulesByStatus(mockSchedules, 'all');
      expect(filtered).toHaveLength(3);
    });

    it('should return empty array when no schedules match', () => {
      const filtered = filterSchedulesByStatus(mockSchedules, 'in_progress');
      expect(filtered).toHaveLength(0);
    });
  });

  describe('calculateVisitDuration', () => {
    it('should calculate duration in minutes', () => {
      const startTime = new Date('2024-01-01T09:00:00Z');
      const endTime = new Date('2024-01-01T10:30:00Z');
      const duration = calculateVisitDuration(startTime, endTime);
      expect(duration).toBe(90);
    });

    it('should return 0 if end time is before start time', () => {
      const startTime = new Date('2024-01-01T10:00:00Z');
      const endTime = new Date('2024-01-01T09:00:00Z');
      const duration = calculateVisitDuration(startTime, endTime);
      expect(duration).toBe(0);
    });
  });

  describe('formatTimestamp', () => {
    it('should format timestamp correctly', () => {
      const timestamp = new Date('2024-01-01T09:30:00Z');
      const formatted = formatTimestamp(timestamp);
      expect(formatted).toBe('Jan 1, 2024 at 9:30 AM');
    });

    it('should handle invalid timestamp', () => {
      const formatted = formatTimestamp(null as any);
      expect(formatted).toBe('Invalid date');
    });
  });

  describe('getStatusColor', () => {
    it('should return correct color for each status', () => {
      expect(getStatusColor('upcoming')).toBe('bg-blue-100 text-blue-800');
      expect(getStatusColor('completed')).toBe('bg-green-100 text-green-800');
      expect(getStatusColor('missed')).toBe('bg-red-100 text-red-800');
      expect(getStatusColor('in_progress')).toBe('bg-yellow-100 text-yellow-800');
    });

    it('should return default color for unknown status', () => {
      expect(getStatusColor('unknown' as any)).toBe('bg-gray-100 text-gray-800');
    });
  });

  describe('validateGeolocation', () => {
    it('should validate valid coordinates', () => {
      const location = { latitude: 40.7128, longitude: -74.0060 };
      const result = validateGeolocation(location);
      expect(result.isValid).toBe(true);
    });

    it('should reject invalid latitude', () => {
      const location = { latitude: 91, longitude: -74.0060 };
      const result = validateGeolocation(location);
      expect(result.isValid).toBe(false);
      expect(result.error).toContain('Latitude must be between -90 and 90');
    });

    it('should reject invalid longitude', () => {
      const location = { latitude: 40.7128, longitude: 181 };
      const result = validateGeolocation(location);
      expect(result.isValid).toBe(false);
      expect(result.error).toContain('Longitude must be between -180 and 180');
    });
  });
});