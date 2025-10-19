import { useState, useEffect } from 'react';
import { useQuery } from '@tanstack/react-query';
import { RefreshCw, Download, Copy, CheckCircle, XCircle, AlertCircle } from 'lucide-react';
import { type Schedule, type Stats } from '../api/standalone';

interface APIResponse {
  endpoint: string;
  method: string;
  status: number;
  success: boolean;
  responseTime: number;
  timestamp: string;
  data?: any;
  error?: string;
}

interface APIResponseViewerProps {
  enabled?: boolean;
  showOnPage?: boolean;
}

export function APIResponseViewer({ enabled = true, showOnPage = false }: APIResponseViewerProps) {
  const [responses, setResponses] = useState<APIResponse[]>([]);
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [lastFetched, setLastFetched] = useState<Date | null>(null);

  // Mock API functions that simulate real API calls
  const mockGetStats = async (): Promise<Stats> => ({
    total: 8,
    missed: 2,
    upcoming: 3,
    completed: 3,
  });

  const mockGetSchedules = async (): Promise<Schedule[]> => [
    {
      id: "1",
      clientName: "John Smith",
      shiftTime: "09:00 - 12:00",
      location: "123 Main St, Anytown",
      status: "upcoming",
      tasks: [
        {
          id: "task1",
          scheduleId: "1",
          name: "Morning Medication",
          description: "Administer morning medication",
          status: "pending",
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    },
  ];

  // Fetch real API data
  const { data: stats, isLoading: statsLoading, error: statsError } = useQuery({
    queryKey: ['real-stats'],
    queryFn: async () => {
      const start = performance.now();
      try {
        // Try real API first
        const response = await fetch('http://localhost:9090/api/v1/schedules/stats');
        if (!response.ok) {
          throw new Error('Failed to fetch stats');
        }
        const data = await response.json();
        const end = performance.now();

        addResponse({
          endpoint: '/api/v1/schedules/stats',
          method: 'GET',
          status: response.status,
          success: true,
          responseTime: Math.round(end - start),
          timestamp: new Date().toISOString(),
          data,
        });

        return data;
      } catch (error) {
        // Fallback to mock data
        const end = performance.now();

        addResponse({
          endpoint: '/api/v1/schedules/stats',
          method: 'GET',
          status: 200, // Mock success
          success: false, // Actually failed but we show mock data
          responseTime: Math.round(end - start),
          timestamp: new Date().toISOString(),
          error: error instanceof Error ? error.message : 'Unknown error',
        });

        return mockGetStats();
      }
    },
    enabled: enabled && autoRefresh,
    refetchInterval: autoRefresh ? 30000 : undefined, // Auto refresh every 30 seconds
  });

  const { data: schedules, isLoading: schedulesLoading, error: schedulesError } = useQuery({
    queryKey: ['real-schedules'],
    queryFn: async () => {
      const start = performance.now();
      try {
        // Try real API first
        const response = await fetch('http://localhost:9090/api/v1/schedules');
        if (!response.ok) {
          throw new Error('Failed to fetch schedules');
        }
        const data = await response.json();
        const end = performance.now();

        addResponse({
          endpoint: '/api/v1/schedules',
          method: 'GET',
          status: response.status,
          success: true,
          responseTime: Math.round(end - start),
          timestamp: new Date().toISOString(),
          data,
        });

        return data;
      } catch (error) {
        // Fallback to mock data
        const end = performance.now();

        addResponse({
          endpoint: '/api/v1/schedules',
          method: 'GET',
          status: 200, // Mock success
          success: false, // Actually failed but we show mock data
          responseTime: Math.round(end - start),
          timestamp: new Date().toISOString(),
          error: error instanceof Error ? error.message : 'Unknown error',
        });

        return mockGetSchedules();
      }
    },
    enabled: enabled && autoRefresh,
    refetchInterval: autoRefresh ? 30000 : undefined,
  });

  const addResponse = (response: APIResponse) => {
    setResponses(prev => {
      const newResponses = [response, ...prev].slice(0, 20); // Keep last 20 responses
      return newResponses;
    });
    setLastFetched(new Date());
  };

  const refreshData = () => {
    setAutoRefresh(false);
    setTimeout(() => setAutoRefresh(true), 100);
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const downloadResponse = (response: APIResponse) => {
    const data = {
      endpoint: response.endpoint,
      method: response.method,
      status: response.status,
      timestamp: response.timestamp,
      responseTime: response.responseTime,
      ...(response.success ? { data: response.data } : { error: response.error }),
    };

    const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = `api-response-${response.timestamp}.json`;
    a.click();
    URL.revokeObjectURL(url);
  };

  if (!enabled) {
    return null;
  }

  return (
    <div className={`${showOnPage ? 'mt-8' : 'fixed top-4 right-4 z-40'} w-full max-w-4xl`}>
      <div className="bg-white rounded-lg shadow-xl border border-gray-200">
        {/* Header */}
        <div className="flex items-center justify-between p-4 border-b border-gray-200">
          <div className="flex items-center space-x-2">
            <div className={`w-2 h-2 rounded-full ${
              statsLoading || schedulesLoading ? 'bg-yellow-500 animate-pulse' : 'bg-green-500'
            }`}></div>
            <h3 className="text-lg font-semibold text-gray-900">
              API Response Monitor
            </h3>
            <span className="text-xs text-gray-500">
              {responses.length} responses
            </span>
          </div>
          <div className="flex items-center space-x-2">
            <button
              onClick={refreshData}
              disabled={statsLoading || schedulesLoading}
              className="flex items-center space-x-1 px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-md disabled:opacity-50"
            >
              <RefreshCw className={`h-3 w-3 ${statsLoading || schedulesLoading ? 'animate-spin' : ''}`} />
              <span>Refresh</span>
            </button>
            <label className="flex items-center space-x-1 text-sm text-gray-600">
              <input
                type="checkbox"
                checked={autoRefresh}
                onChange={(e) => setAutoRefresh(e.target.checked)}
                className="rounded border-gray-300"
              />
              <span>Auto Refresh</span>
            </label>
          </div>
        </div>

        {/* API Status Cards */}
        <div className="p-4 grid grid-cols-1 md:grid-cols-2 gap-4 border-b border-gray-200">
          <div className="p-3 bg-blue-50 rounded-lg border border-blue-200">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium text-blue-900">Stats API</h4>
                <p className="text-sm text-blue-700">
                  /api/v1/schedules/stats
                </p>
              </div>
              {statsError ? (
                <XCircle className="h-5 w-5 text-red-500" />
              ) : statsLoading ? (
                <div className="h-5 w-5 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
              ) : (
                <CheckCircle className="h-5 w-5 text-green-500" />
              )}
            </div>
            {statsError && (
              <p className="text-xs text-red-600 mt-1">{statsError.message}</p>
            )}
          </div>

          <div className="p-3 bg-green-50 rounded-lg border border-green-200">
            <div className="flex items-center justify-between">
              <div>
                <h4 className="font-medium text-green-900">Schedules API</h4>
                <p className="text-sm text-green-700">
                  /api/v1/schedules
                </p>
              </div>
              {schedulesError ? (
                <XCircle className="h-5 w-5 text-red-500" />
              ) : schedulesLoading ? (
                <div className="h-5 w-5 border-2 border-green-500 border-t-transparent rounded-full animate-spin"></div>
              ) : (
                <CheckCircle className="h-5 w-5 text-green-500" />
              )}
            </div>
            {schedulesError && (
              <p className="text-xs text-red-600 mt-1">{schedulesError.message}</p>
            )}
          </div>
        </div>

        {/* Response History */}
        <div className="p-4 max-h-64 overflow-y-auto">
          <h4 className="font-medium text-gray-900 mb-3">Recent Responses</h4>
          {responses.length === 0 ? (
            <div className="text-center py-8 text-gray-500">
              <AlertCircle className="h-8 w-8 mx-auto mb-2 opacity-50" />
              <p className="text-sm">No API responses recorded yet</p>
            </div>
          ) : (
            <div className="space-y-2">
              {responses.map((response, index) => (
                <div
                  key={index}
                  className={`p-3 rounded-lg border ${
                    response.success
                      ? 'bg-green-50 border-green-200'
                      : 'bg-red-50 border-red-200'
                  }`}
                >
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center space-x-2">
                      <span className="font-mono text-xs font-bold text-gray-600">
                        {response.method}
                      </span>
                      <span className="font-mono text-xs bg-gray-100 px-1 rounded">
                        {response.endpoint}
                      </span>
                      <span className={`text-xs px-2 py-0.5 rounded-full font-medium ${
                        response.success ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                      }`}>
                        {response.status}
                      </span>
                    </div>
                    <span className="text-xs font-mono text-gray-500">
                      {response.responseTime}ms
                    </span>
                  </div>

                  {response.error && (
                    <div className="text-xs text-red-600 flex items-start space-x-1">
                      <XCircle className="h-3 w-3 mt-0.5 flex-shrink-0" />
                      <span>{response.error}</span>
                    </div>
                  )}

                  <div className="flex justify-between items-center mt-2">
                    <span className="text-xs text-gray-500">
                      {new Date(response.timestamp).toLocaleTimeString()}
                    </span>
                    <div className="flex space-x-1">
                      <button
                        onClick={() => copyToClipboard(JSON.stringify(response.data, null, 2))}
                        className="p-1 hover:bg-gray-100 rounded text-gray-600"
                        title="Copy response"
                      >
                        <Copy className="h-3 w-3" />
                      </button>
                      <button
                        onClick={() => downloadResponse(response)}
                        className="p-1 hover:bg-gray-100 rounded text-gray-600"
                        title="Download response"
                      >
                        <Download className="h-3 w-3" />
                      </button>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="p-3 bg-gray-50 border-t border-gray-200 text-xs">
          <div className="flex justify-between items-center">
            <span className="text-gray-600">
              {lastFetched ? `Last updated: ${lastFetched.toLocaleTimeString()}` : 'Waiting for data...'}
            </span>
            {autoRefresh && (
              <span className="text-green-600">Auto-refreshing every 30s</span>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}