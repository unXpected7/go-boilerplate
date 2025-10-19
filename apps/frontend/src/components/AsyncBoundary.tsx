import { Component, ErrorInfo, ReactNode, Suspense } from 'react';
import { ErrorBoundary } from './ErrorBoundary';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class AsyncBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      hasError: false,
      error: null,
    };
  }

  static getDerivedStateFromError(error: Error): State {
    return {
      hasError: true,
      error,
    };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    this.setState({
      error,
    });

    // Log the error
    console.error('Async component error:', error, errorInfo);
  }

  render() {
    const { fallback, children } = this.props;

    if (this.state.hasError) {
      return fallback || (
        <div className="flex items-center justify-center h-64">
          <div className="text-center">
            <p className="text-red-600 mb-4">Failed to load component</p>
            <button
              onClick={() => this.setState({ hasError: false, error: null })}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Retry
            </button>
          </div>
        </div>
      );
    }

    return (
      <ErrorBoundary>
        <Suspense fallback={fallback || this.getDefaultFallback()}>
          {children}
        </Suspense>
      </ErrorBoundary>
    );
  }

  private getDefaultFallback = () => {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="text-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-2"></div>
          <p className="text-gray-600">Loading...</p>
        </div>
      </div>
    );
  };
}