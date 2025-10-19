import { useEffect, useState } from 'react';

interface ServiceWorkerStatus {
  isSupported: boolean;
  isRegistered: boolean;
  isReady: boolean;
  error: string | null;
}

export function useServiceWorker(): ServiceWorkerStatus & {
  update: () => Promise<void>;
  skipWaiting: () => Promise<void>;
} {
  const [status, setStatus] = useState<ServiceWorkerStatus>({
    isSupported: false,
    isRegistered: false,
    isReady: false,
    error: null,
  });

  useEffect(() => {
    // Check if Service Worker is supported
    if ('serviceWorker' in navigator) {
      setStatus(prev => ({ ...prev, isSupported: true }));

      // Register Service Worker
      navigator.serviceWorker
        .register('/sw.js')
        .then((registration) => {
          console.log('Service Worker registered:', registration);
          setStatus(prev => ({ ...prev, isRegistered: true }));

          // Check if service worker is ready
          if (registration.active) {
            setStatus(prev => ({ ...prev, isReady: true }));
          } else {
            registration.addEventListener('updatefound', () => {
              const installingWorker = registration.installing;
              installingWorker?.addEventListener('statechange', () => {
                if (installingWorker.state === 'activated') {
                  setStatus(prev => ({ ...prev, isReady: true }));
                }
              });
            });
          }
        })
        .catch((error) => {
          console.error('Service Worker registration failed:', error);
          setStatus(prev => ({
            ...prev,
            error: error.message || 'Service Worker registration failed',
          }));
        });

      // Listen for service worker messages
      navigator.serviceWorker.addEventListener('message', (event) => {
        console.log('Service Worker message:', event.data);
      });

      // Handle controller change
      navigator.serviceWorker.addEventListener('controllerchange', () => {
        console.log('Service Worker controller changed');
        window.location.reload();
      });

      // Cleanup
      return () => {
        navigator.serviceWorker.removeEventListener('message', () => {});
        navigator.serviceWorker.removeEventListener('controllerchange', () => {});
      };
    } else {
      setStatus(prev => ({
        ...prev,
        error: 'Service Workers are not supported in this browser',
      }));
    }
  }, []);

  const update = async (): Promise<void> => {
    if (!status.isRegistered) return;

    try {
      const registration = await navigator.serviceWorker.getRegistration();
      if (registration) {
        await registration.update();
        console.log('Service Worker update initiated');
      }
    } catch (error) {
      console.error('Service Worker update failed:', error);
      throw error;
    }
  };

  const skipWaiting = async (): Promise<void> => {
    if (!status.isRegistered) return;

    try {
      const registration = await navigator.serviceWorker.getRegistration();
      if (registration) {
        await registration.waiting?.postMessage({ type: 'SKIP_WAITING' });
      }
    } catch (error) {
      console.error('Skip waiting failed:', error);
      throw error;
    }
  };

  return {
    ...status,
    update,
    skipWaiting,
  };
}

export function useBackgroundSync() {
  const [isSyncing, setIsSyncing] = useState(false);
  const [lastSync, setLastSync] = useState<Date | null>(null);

  const registerSync = async (tag: string): Promise<void> => {
    if (!('serviceWorker' in navigator) || !('sync' in window.ServiceWorkerRegistration.prototype)) {
      console.warn('Background Sync is not supported');
      return;
    }

    try {
      setIsSyncing(true);
      const registration = await navigator.serviceWorker.ready;
      await registration.sync.register(tag);
      console.log(`Sync registered for tag: ${tag}`);
    } catch (error) {
      console.error('Background Sync registration failed:', error);
      throw error;
    } finally {
      setIsSyncing(false);
      setLastSync(new Date());
    }
  };

  return {
    isSyncing,
    lastSync,
    registerSync,
  };
}