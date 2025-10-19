import { useState, useEffect, useCallback } from "react";
import { type Location } from "../api/types";

interface GeolocationState {
  location: Location | null;
  error: string | null;
  isLoading: boolean;
  isSupported: boolean;
}

export function useGeolocation(options: PositionOptions = {}) {
  const [state, setState] = useState<GeolocationState>({
    location: null,
    error: null,
    isLoading: false,
    isSupported: typeof navigator !== "undefined" && "geolocation" in navigator,
  });

  const getLocation = useCallback(() => {
    if (!state.isSupported) {
      setState(prev => ({
        ...prev,
        error: "Geolocation is not supported by this browser.",
        isLoading: false,
      }));
      return;
    }

    setState(prev => ({ ...prev, isLoading: true, error: null }));

    navigator.geolocation.getCurrentPosition(
      (position) => {
        const location: Location = {
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy,
          timestamp: new Date(),
        };
        setState({
          location,
          error: null,
          isLoading: false,
          isSupported: true,
        });
      },
      (error) => {
        let errorMessage = "Failed to get location.";

        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = "Location access denied by user.";
            break;
          case error.POSITION_UNAVAILABLE:
            errorMessage = "Location information is unavailable.";
            break;
          case error.TIMEOUT:
            errorMessage = "Location request timed out.";
            break;
          default:
            errorMessage = "An unknown error occurred while retrieving location.";
            break;
        }

        setState({
          location: null,
          error: errorMessage,
          isLoading: false,
          isSupported: true,
        });
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0,
        ...options,
      }
    );
  }, [options, state.isSupported]);

  const watchLocation = useCallback(() => {
    if (!state.isSupported) return null;

    const watchId = navigator.geolocation.watchPosition(
      (position) => {
        const location: Location = {
          latitude: position.coords.latitude,
          longitude: position.coords.longitude,
          accuracy: position.coords.accuracy,
          timestamp: new Date(),
        };
        setState(prev => ({
          ...prev,
          location,
          error: null,
          isLoading: false,
        }));
      },
      (error) => {
        let errorMessage = "Failed to watch location.";

        switch (error.code) {
          case error.PERMISSION_DENIED:
            errorMessage = "Location access denied by user.";
            break;
          case error.POSITION_UNAVAILABLE:
            errorMessage = "Location information is unavailable.";
            break;
          case error.TIMEOUT:
            errorMessage = "Location request timed out.";
            break;
          default:
            errorMessage = "An unknown error occurred while watching location.";
            break;
        }

        setState(prev => ({
          ...prev,
          error: errorMessage,
          isLoading: false,
        }));
      },
      {
        enableHighAccuracy: true,
        timeout: 10000,
        maximumAge: 0,
        ...options,
      }
    );

    return watchId;
  }, [options, state.isSupported]);

  const clearWatch = useCallback((watchId: number) => {
    if (navigator.geolocation) {
      navigator.geolocation.clearWatch(watchId);
    }
  }, []);

  useEffect(() => {
    return () => {
      // Cleanup any active watches when component unmounts
      if (navigator.geolocation) {
        navigator.geolocation.clearWatch(0);
      }
    };
  }, []);

  return {
    ...state,
    getLocation,
    watchLocation,
    clearWatch,
  };
}