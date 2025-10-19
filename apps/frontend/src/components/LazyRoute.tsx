import { lazy, Suspense, ComponentType } from "react";
import { Loader2 } from "lucide-react";

interface LazyRouteProps {
  children: ComponentType<any>;
  fallback?: React.ReactNode;
}

export function LazyRoute({ children: Component, fallback }: LazyRouteProps) {
  const FallbackComponent = fallback || (
    <div className="flex items-center justify-center h-64">
      <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      <span className="ml-2 text-gray-600">Loading...</span>
    </div>
  );

  return (
    <Suspense fallback={FallbackComponent}>
      {typeof Component === 'function' ? <Component /> : Component}
    </Suspense>
  );
}

// Lazy load components for better performance
export const LazyHome = lazy(() => import("../pages/Home"));
export const LazyScheduleDetails = lazy(() => import("../pages/ScheduleDetails"));