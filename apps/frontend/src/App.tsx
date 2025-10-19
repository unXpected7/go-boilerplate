import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { Toaster } from "sonner";
import { Layout } from "./components/Layout";
import { LazyRoute, LazyHome, LazyScheduleDetails } from "./components/LazyRoute";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { OfflineBanner, OfflineIndicator } from "./components/OfflineBanner";
import { SkipToContent } from "./components/SkipToContent";
import { AccessibilityControls } from "./components/AccessibilityControls";
import { useAccessibilityShortcuts } from "./hooks/useKeyboardShortcut";

const queryClient = new QueryClient();

function App() {
  // Initialize accessibility shortcuts
  useAccessibilityShortcuts();

  return (
    <ErrorBoundary>
      <QueryClientProvider client={queryClient}>
        <Router>
          <div className="min-h-screen bg-gray-50">
            <SkipToContent mainContentId="main-content" />
            <OfflineBanner />
            <Layout>
              <Routes>
                <Route
                  path="/"
                  element={
                    <LazyRoute>
                      <LazyHome />
                    </LazyRoute>
                  }
                />
                <Route
                  path="/schedules/:id"
                  element={
                    <LazyRoute>
                      <LazyScheduleDetails />
                    </LazyRoute>
                  }
                />
              </Routes>
            </Layout>
            <OfflineIndicator />
            <AccessibilityControls />
          </div>
          <Toaster />
        </Router>
      </QueryClientProvider>
    </ErrorBoundary>
  );
}

export default App;
