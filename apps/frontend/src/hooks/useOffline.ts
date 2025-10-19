import { useState, useEffect } from 'react';

export interface OfflineStatus {
  isOnline: boolean;
  isOffline: boolean;
  lastOnline: Date | null;
  lastOffline: Date | null;
  connectionType?: string;
  effectiveType?: string;
  downlink?: number;
  rtt?: number;
}

export function useOffline(): OfflineStatus {
  const [status, setStatus] = useState<OfflineStatus>({
    isOnline: true,
    isOffline: false,
    lastOnline: null,
    lastOffline: null,
    connectionType: 'unknown',
    effectiveType: 'unknown',
    downlink: 0,
    rtt: 0,
  });

  useEffect(() => {
    const handleOnline = () => {
      setStatus(prev => ({
        ...prev,
        isOnline: true,
        isOffline: false,
        lastOnline: new Date(),
        lastOffline: prev.isOffline ? new Date() : prev.lastOffline,
      }));
    };

    const handleOffline = () => {
      setStatus(prev => ({
        ...prev,
        isOnline: false,
        isOffline: true,
        lastOffline: new Date(),
        lastOnline: prev.isOnline ? new Date() : prev.lastOnline,
      }));
    };

    const handleConnectionChange = () => {
      const connection = (navigator as any).connection ||
                        (navigator as any).mozConnection ||
                        (navigator as any).webkitConnection;

      if (connection) {
        setStatus(prev => ({
          ...prev,
          connectionType: connection.type || 'unknown',
          effectiveType: connection.effectiveType || 'unknown',
          downlink: connection.downlink || 0,
          rtt: connection.rtt || 0,
        }));
      }
    };

    // Set up event listeners
    window.addEventListener('online', handleOnline);
    window.addEventListener('offline', handleOffline);
    window.addEventListener('connectionchange', handleConnectionChange);

    // Initial check
    handleConnectionChange();

    // Check initial online status
    if (!navigator.onLine) {
      handleOffline();
    }

    return () => {
      window.removeEventListener('online', handleOnline);
      window.removeEventListener('offline', handleOffline);
      window.removeEventListener('connectionchange', handleConnectionChange);
    };
  }, []);

  return status;
}

export function useNetworkAwareRequest(requestFn: () => Promise<any>) {
  const { isOnline } = useOffline();
  const [data, setData] = useState<any>(null);
  const [error, setError] = useState<Error | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const execute = async () => {
    if (!isOnline) {
      setError(new Error('Network is offline. Please check your connection.'));
      return;
    }

    setIsLoading(true);
    setError(null);

    try {
      const result = await requestFn();
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err : new Error('Request failed'));
    } finally {
      setIsLoading(false);
    }
  };

  return {
    data,
    error,
    isLoading,
    isOnline,
    execute,
    refetch: execute,
  };
}