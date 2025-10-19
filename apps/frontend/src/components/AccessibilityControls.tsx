import { useState, useCallback, useEffect } from 'react';
import {
  Sun,
  Moon,
  Volume2,
  VolumeX,
  Eye,
  EyeOff,
  Type,
  Ruler,
  ZoomIn,
  ZoomOut
} from 'lucide-react';

interface AccessibilitySettings {
  highContrast: boolean;
  reducedMotion: boolean;
  screenReader: boolean;
  fontSize: 'small' | 'medium' | 'large' | 'x-large';
  letterSpacing: 'normal' | 'wide' | 'wider';
  lineHeight: 'normal' | 'relaxed' | 'loose';
}

const DEFAULT_SETTINGS: AccessibilitySettings = {
  highContrast: false,
  reducedMotion: false,
  screenReader: false,
  fontSize: 'medium',
  letterSpacing: 'normal',
  lineHeight: 'normal',
};

export function AccessibilityControls() {
  const [isOpen, setIsOpen] = useState(false);
  const [settings, setSettings] = useState<AccessibilitySettings>(DEFAULT_SETTINGS);

  const applyStyles = useCallback(() => {
    const root = document.documentElement;

    // Apply high contrast mode
    if (settings.highContrast) {
      root.classList.add('high-contrast');
    } else {
      root.classList.remove('high-contrast');
    }

    // Apply reduced motion
    if (settings.reducedMotion) {
      root.classList.add('reduce-motion');
    } else {
      root.classList.remove('reduce-motion');
    }

    // Apply font size
    root.style.fontSize = {
      small: '14px',
      medium: '16px',
      large: '18px',
      'x-large': '20px',
    }[settings.fontSize];

    // Apply letter spacing
    root.style.letterSpacing = {
      normal: 'normal',
      wide: '0.05em',
      wider: '0.1em',
    }[settings.letterSpacing];

    // Apply line height
    root.style.lineHeight = {
      normal: '1.5',
      relaxed: '1.625',
      loose: '2',
    }[settings.lineHeight];

    // Announce changes to screen readers
    if (settings.screenReader) {
      const announcement = document.createElement('div');
      announcement.setAttribute('role', 'status');
      announcement.setAttribute('aria-live', 'polite');
      announcement.className = 'sr-only';
      announcement.textContent = 'Accessibility settings updated';
      document.body.appendChild(announcement);
      setTimeout(() => document.body.removeChild(announcement), 1000);
    }
  }, [settings]);

  useEffect(() => {
    applyStyles();
  }, [applyStyles]);

  const toggleSetting = (key: keyof AccessibilitySettings) => {
    setSettings(prev => ({ ...prev, [key]: !prev[key] }));
  };

  const updateSetting = <K extends keyof AccessibilitySettings>(
    key: K,
    value: AccessibilitySettings[K]
  ) => {
    setSettings(prev => ({ ...prev, [key]: value }));
  };

  return (
    <div className="fixed bottom-4 left-4 z-40">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="bg-blue-600 text-white p-3 rounded-full shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        aria-label="Accessibility controls"
        aria-expanded={isOpen}
        aria-controls="accessibility-panel"
      >
        <Eye className="h-5 w-5" />
      </button>

      {isOpen && (
        <div
          id="accessibility-panel"
          className="absolute bottom-16 left-0 bg-white rounded-lg shadow-xl p-6 w-80 max-w-sm"
          role="dialog"
          aria-modal="true"
          aria-labelledby="accessibility-title"
        >
          <h2
            id="accessibility-title"
            className="text-lg font-semibold text-gray-900 mb-4"
          >
            Accessibility Settings
          </h2>

          <div className="space-y-4">
            {/* High Contrast */}
            <div className="flex items-center justify-between">
              <label className="flex items-center space-x-3">
                <Sun className="h-4 w-4 text-gray-600" />
                <span className="text-sm font-medium text-gray-700">
                  High Contrast
                </span>
              </label>
              <button
                onClick={() => toggleSetting('highContrast')}
                className={`relative inline-flex h-6 w-11 items-center rounded-full ${
                  settings.highContrast ? 'bg-blue-600' : 'bg-gray-200'
                }`}
                aria-pressed={settings.highContrast}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition ${
                    settings.highContrast ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Reduced Motion */}
            <div className="flex items-center justify-between">
              <label className="flex items-center space-x-3">
                <Type className="h-4 w-4 text-gray-600" />
                <span className="text-sm font-medium text-gray-700">
                  Reduced Motion
                </span>
              </label>
              <button
                onClick={() => toggleSetting('reducedMotion')}
                className={`relative inline-flex h-6 w-11 items-center rounded-full ${
                  settings.reducedMotion ? 'bg-blue-600' : 'bg-gray-200'
                }`}
                aria-pressed={settings.reducedMotion}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition ${
                    settings.reducedMotion ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Screen Reader Mode */}
            <div className="flex items-center justify-between">
              <label className="flex items-center space-x-3">
                <Volume2 className="h-4 w-4 text-gray-600" />
                <span className="text-sm font-medium text-gray-700">
                  Screen Reader Mode
                </span>
              </label>
              <button
                onClick={() => toggleSetting('screenReader')}
                className={`relative inline-flex h-6 w-11 items-center rounded-full ${
                  settings.screenReader ? 'bg-blue-600' : 'bg-gray-200'
                }`}
                aria-pressed={settings.screenReader}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition ${
                    settings.screenReader ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>

            {/* Font Size */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700">
                Font Size
              </label>
              <div className="flex space-x-2">
                {(['small', 'medium', 'large', 'x-large'] as const).map((size) => (
                  <button
                    key={size}
                    onClick={() => updateSetting('fontSize', size)}
                    className={`flex-1 py-2 px-3 text-xs rounded-md border ${
                      settings.fontSize === size
                        ? 'bg-blue-600 text-white border-blue-600'
                        : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                    }`}
                    aria-pressed={settings.fontSize === size}
                  >
                    {size.charAt(0).toUpperCase() + size.slice(1)}
                  </button>
                ))}
              </div>
            </div>

            {/* Letter Spacing */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700">
                Letter Spacing
              </label>
              <div className="flex space-x-2">
                {(['normal', 'wide', 'wider'] as const).map((spacing) => (
                  <button
                    key={spacing}
                    onClick={() => updateSetting('letterSpacing', spacing)}
                    className={`flex-1 py-2 px-3 text-xs rounded-md border ${
                      settings.letterSpacing === spacing
                        ? 'bg-blue-600 text-white border-blue-600'
                        : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                    }`}
                    aria-pressed={settings.letterSpacing === spacing}
                  >
                    {spacing.charAt(0).toUpperCase() + spacing.slice(1)}
                  </button>
                ))}
              </div>
            </div>

            {/* Line Height */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-gray-700">
                Line Height
              </label>
              <div className="flex space-x-2">
                {(['normal', 'relaxed', 'loose'] as const).map((height) => (
                  <button
                    key={height}
                    onClick={() => updateSetting('lineHeight', height)}
                    className={`flex-1 py-2 px-3 text-xs rounded-md border ${
                      settings.lineHeight === height
                        ? 'bg-blue-600 text-white border-blue-600'
                        : 'bg-white text-gray-700 border-gray-300 hover:bg-gray-50'
                    }`}
                    aria-pressed={settings.lineHeight === height}
                  >
                    {height.charAt(0).toUpperCase() + height.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          </div>

          <button
            onClick={() => {
              setSettings(DEFAULT_SETTINGS);
              setIsOpen(false);
            }}
            className="mt-4 w-full py-2 px-4 bg-gray-100 text-gray-700 rounded-md hover:bg-gray-200 text-sm font-medium"
          >
            Reset to Defaults
          </button>
        </div>
      )}
    </div>
  );
}