import { useOffline } from '../hooks/useOffline';
import { Wifi, WifiOff, RefreshCw } from 'lucide-react';
import { useEffect, useState } from 'react';

export function OfflineBanner() {
  const { isOnline, isOffline, lastOnline, lastOffline } = useOffline();
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Show banner when offline
    if (isOffline) {
      setIsVisible(true);
    } else {
      // Hide after a delay when coming back online
      const timer = setTimeout(() => {
        setIsVisible(false);
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [isOffline]);

  if (!isVisible) return null;

  const formatDate = (date: Date | null) => {
    if (!date) return '';
    return date.toLocaleTimeString();
  };

  return (
    <div className={`fixed top-0 left-0 right-0 z-50 p-4 transition-all duration-300 ${
      isOnline
        ? 'bg-green-100 border-b border-green-300'
        : 'bg-red-100 border-b border-red-300'
    }`}>
      <div className="max-w-7xl mx-auto flex items-center justify-between">
        <div className="flex items-center space-x-3">
          {isOffline ? (
            <WifiOff className="h-5 w-5 text-red-600" />
          ) : (
            <Wifi className="h-5 w-5 text-green-600" />
          )}

          <div>
            <p className={`text-sm font-medium ${
              isOnline ? 'text-green-800' : 'text-red-800'
            }`}>
              {isOnline
                ? `✓ Back online since ${formatDate(lastOnline)}`
                : `✗ Offline since ${formatDate(lastOffline)}`
              }
            </p>
            {isOffline && (
              <p className="text-xs text-red-600">
                Some features may be limited. Changes will be synced when back online.
              </p>
            )}
          </div>
        </div>

        {isOffline && (
          <button
            onClick={() => window.location.reload()}
            className="inline-flex items-center px-3 py-1 border border-red-300 text-xs font-medium rounded text-red-700 bg-white hover:bg-red-50"
          >
            <RefreshCw className="h-3 w-3 mr-1" />
            Retry
          </button>
        )}
      </div>
    </div>
  );
}

export function OfflineIndicator() {
  const { isOnline, connectionType, effectiveType, downlink } = useOffline();

  return (
    <div className="fixed bottom-4 right-4 z-40">
      <div className="bg-white rounded-lg shadow-lg p-3 flex items-center space-x-2">
        {isOnline ? (
          <Wifi className="h-4 w-4 text-green-500" />
        ) : (
          <WifiOff className="h-4 w-4 text-red-500" />
        )}

        <span className="text-xs font-medium text-gray-700">
          {isOnline ? 'Online' : 'Offline'}
        </span>

        {isOnline && connectionType && (
          <span className="text-xs text-gray-500">
            ({effectiveType} • {downlink}Mbps)
          </span>
        )}
      </div>
    </div>
  );
}