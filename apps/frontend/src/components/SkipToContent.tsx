import { useRef } from 'react';
import { useKeyboardShortcut } from '../hooks/useKeyboardShortcut';

interface SkipToContentProps {
  mainContentId: string;
}

export function SkipToContent({ mainContentId }: SkipToContentProps) {
  const mainContentRef = useRef<HTMLElement>(null);

  useKeyboardShortcut(';', () => {
    mainContentRef.current?.focus();
  });

  const handleSkip = () => {
    mainContentRef.current?.focus();
  };

  return (
    <a
      href={`#${mainContentId}`}
      onClick={handleSkip}
      className="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:bg-blue-600 focus:text-white focus:p-2 focus:rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
    >
      Skip to main content
    </a>
  );
}