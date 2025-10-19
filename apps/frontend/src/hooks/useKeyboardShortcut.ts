import { useEffect, useCallback } from 'react';

interface ShortcutOptions {
  preventDefault?: boolean;
  stopPropagation?: boolean;
  element?: HTMLElement | Window;
}

export function useKeyboardShortcut(
  key: string,
  callback: () => void,
  options: ShortcutOptions = {}
): void {
  const { preventDefault = true, stopPropagation = false, element = window } = options;

  const handleKeyPress = useCallback((event: KeyboardEvent) => {
    // Check if the pressed key matches the shortcut key
    if (event.key === key) {
      if (preventDefault) {
        event.preventDefault();
      }
      if (stopPropagation) {
        event.stopPropagation();
      }
      callback();
    }
  }, [key, callback, preventDefault, stopPropagation]);

  useEffect(() => {
    element.addEventListener('keydown', handleKeyPress);

    return () => {
      element.removeEventListener('keydown', handleKeyPress);
    };
  }, [handleKeyPress, element]);
}

export function useAccessibilityShortcuts() {
  useKeyboardShortcut('Escape', () => {
    // Close modals, dropdowns, etc.
    document.activeElement instanceof HTMLElement && document.activeElement.blur();
  });

  useKeyboardShortcut('Tab', () => {
    // Enhanced tab navigation could be added here
    console.log('Tab key pressed for accessibility');
  });

  useKeyboardShortcut('Shift+Tab', () => {
    // Shift+Tab navigation
    console.log('Shift+Tab key pressed for accessibility');
  });
}