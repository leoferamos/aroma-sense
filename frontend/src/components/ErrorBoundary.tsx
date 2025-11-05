import React from 'react';

type ErrorBoundaryState = { hasError: boolean };

export default class ErrorBoundary extends React.Component<React.PropsWithChildren, ErrorBoundaryState> {
  constructor(props: React.PropsWithChildren) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(): ErrorBoundaryState {
    return { hasError: true };
  }

  componentDidCatch(error: unknown, errorInfo: unknown) {
    console.error('Uncaught error in component tree:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="bg-white shadow rounded-lg p-8 text-center max-w-md">
            <h1 className="text-xl font-semibold text-gray-900">Something went wrong.</h1>
            <p className="mt-2 text-gray-600">Please refresh the page or try again later.</p>
            <div className="mt-6 flex flex-col gap-3">
              <button
                type="button"
                onClick={() => this.setState({ hasError: false })}
                className="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors"
              >
                Try Again
              </button>
              <button
                type="button"
                onClick={() => { window.location.assign('/products'); }}
                className="px-4 py-2 rounded-lg border border-gray-300 text-gray-800 hover:bg-gray-50 transition-colors"
              >
                Go to Products
              </button>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
