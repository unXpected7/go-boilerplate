import { useCallback } from 'react';
import { toast } from 'sonner';

export interface ErrorHandlerOptions {
  logError?: boolean;
  showUserMessage?: boolean;
  userMessage?: string;
  retry?: boolean;
  maxRetries?: number;
  delay?: number;
}

export function useErrorHandler() {
  const handleError = useCallback((
    error: Error | unknown,
    options: ErrorHandlerOptions = {}
  ) => {
    const {
      logError = true,
      showUserMessage = true,
      userMessage,
      retry = false,
      maxRetries = 3,
      delay = 1000,
    } = options;

    // Convert unknown error to Error
    const errorObj = error instanceof Error ? error : new Error(String(error));

    // Log error to console
    if (logError) {
      console.error('Error handled by useErrorHandler:', errorObj);
    }

    // Show user message
    if (showUserMessage) {
      const message = userMessage || errorObj.message || 'An unexpected error occurred';
      toast.error(message);
    }

    // Handle retry logic
    if (retry) {
      let attempt = 0;

      const attemptRetry = async (): Promise<boolean> => {
        attempt++;

        if (attempt >= maxRetries) {
          toast.error('Operation failed after multiple attempts. Please try again later.');
          return false;
        }

        await new Promise(resolve => setTimeout(resolve, delay * attempt));

        try {
          // This would be implemented based on the specific use case
          // For now, just simulate a retry
          console.log(`Retry attempt ${attempt} of ${maxRetries}`);
          return true;
        } catch (retryError) {
          if (attempt >= maxRetries - 1) {
            throw retryError;
          }
          return attemptRetry();
        }
      };

      return attemptRetry();
    }

    return false;
  }, []);

  return { handleError };
}