import { useState, useCallback, useEffect } from 'react';

interface RetryOptions {
  maxRetries?: number;
  delay?: number;
  backoff?: boolean;
}

export function useRetry<T>(options: RetryOptions = {}) {
  const {
    maxRetries = 3,
    delay = 1000,
    backoff = true,
  } = options;

  const [retryCount, setRetryCount] = useState(0);
  const [isRetrying, setIsRetrying] = useState(false);

  const executeWithRetry = useCallback(async (
    fn: () => Promise<T>,
    onSuccess?: (data: T) => void,
    onError?: (error: Error) => void
  ): Promise<boolean> => {
    setIsRetrying(true);

    try {
      const data = await fn();
      onSuccess?.(data);
      setRetryCount(0);
      return true;
    } catch (error) {
      const err = error instanceof Error ? error : new Error(String(error));

      if (retryCount < maxRetries - 1) {
        const retryDelay = backoff ? delay * Math.pow(2, retryCount) : delay;

        await new Promise(resolve => setTimeout(resolve, retryDelay));

        setRetryCount(prev => prev + 1);
        const retrySuccess = await executeWithRetry(fn, onSuccess, onError);
        return retrySuccess;
      } else {
        onError?.(err);
        return false;
      }
    } finally {
      setIsRetrying(false);
    }
  }, [retryCount, maxRetries, delay, backoff]);

  const resetRetry = useCallback(() => {
    setRetryCount(0);
    setIsRetrying(false);
  }, []);

  const canRetry = retryCount < maxRetries;

  return {
    executeWithRetry,
    resetRetry,
    canRetry,
    isRetrying,
    retryCount,
  };
}