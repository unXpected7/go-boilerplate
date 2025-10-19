import { renderHook, act } from '@testing-library/react';
import { useOffline, useNetworkAwareRequest } from '../../hooks/useOffline';

// Mock navigator.onLine
Object.defineProperty(navigator, 'onLine', {
  writable: true,
  value: true,
});

// Mock connection API
Object.defineProperty(navigator, 'connection', {
  writable: true,
  value: {
    type: 'wifi',
    effectiveType: '4g',
    downlink: 10,
    rtt: 50,
  },
});

describe('useOffline Hook', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should initially show as online', () => {
    const { result } = renderHook(() => useOffline());

    expect(result.current.isOnline).toBe(true);
    expect(result.current.isOffline).toBe(false);
  });

  it('should detect when going offline', () => {
    const { result } = renderHook(() => useOffline());

    // Simulate going offline
    act(() => {
      window.dispatchEvent(new Event('offline'));
    });

    expect(result.current.isOnline).toBe(false);
    expect(result.current.isOffline).toBe(true);
    expect(result.current.lastOffline).not.toBeNull();
  });

  it('should detect when coming back online', () => {
    const { result } = renderHook(() => useOffline());

    // Simulate going offline and back online
    act(() => {
      window.dispatchEvent(new Event('offline'));
      window.dispatchEvent(new Event('online'));
    });

    expect(result.current.isOnline).toBe(true);
    expect(result.current.isOffline).toBe(false);
    expect(result.current.lastOnline).not.toBeNull();
  });

  it('should handle connection information', () => {
    const { result } = renderHook(() => useOffline());

    expect(result.current.connectionType).toBe('wifi');
    expect(result.current.effectiveType).toBe('4g');
    expect(result.current.downlink).toBe(10);
    expect(result.current.rtt).toBe(50);
  });
});

describe('useNetworkAwareRequest Hook', () => {
  const mockRequestFn = jest.fn()
    .mockResolvedValue({ success: true })
    .mockName('mockRequestFn');

  beforeEach(() => {
    mockRequestFn.mockClear();
  });

  it('should execute request successfully when online', async () => {
    const { result } = renderHook(() => useNetworkAwareRequest(mockRequestFn));

    await act(async () => {
      await result.current.execute();
    });

    expect(result.current.data).toEqual({ success: true });
    expect(result.current.error).toBeNull();
    expect(result.current.isOnline).toBe(true);
  });

  it('should not execute request when offline', async () => {
    const { result } = renderHook(() => useNetworkAwareRequest(mockRequestFn));

    // Simulate being offline
    act(() => {
      window.dispatchEvent(new Event('offline'));
    });

    await act(async () => {
      await result.current.execute();
    });

    expect(result.current.error).toEqual(expect.any(Error));
    expect(mockRequestFn).not.toHaveBeenCalled();
  });

  it('should handle request failure', async () => {
    const errorRequestFn = jest.fn().mockRejectedValue(new Error('Request failed'));
    const { result } = renderHook(() => useNetworkAwareRequest(errorRequestFn));

    await act(async () => {
      await result.current.execute();
    });

    expect(result.current.error).toEqual(expect.any(Error));
    expect(result.current.data).toBeUndefined();
  });
});