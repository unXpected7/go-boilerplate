import { useState } from 'react';
import { Copy, ExternalLink, Play, RefreshCw, AlertCircle, CheckCircle } from 'lucide-react';

interface APIEndpoint {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  path: string;
  description: string;
  parameters?: string[];
  example?: any;
}

interface APIResponse {
  status: number;
  data?: any;
  error?: string;
  responseTime: number;
  timestamp: string;
}

interface APIDocsTesterProps {
  endpoints?: APIEndpoint[];
}

export function APIDocsTester({ endpoints = [] }: APIDocsTesterProps) {
  const [selectedEndpoint, setSelectedEndpoint] = useState<APIEndpoint | null>(null);
  const [response, setResponse] = useState<APIResponse | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [customParams, setCustomParams] = useState('{}');

  // Default API endpoints
  const defaultEndpoints: APIEndpoint[] = [
    {
      method: 'GET',
      path: '/status',
      description: 'Get system health status',
      parameters: [],
      example: {
        status: 'healthy',
        environment: 'development',
        checks: {
          database: { status: 'healthy', responseTime: '219µs' },
          redis: { status: 'healthy', responseTime: '311µs' },
        },
      },
    },
    {
      method: 'GET',
      path: '/api/v1/schedules',
      description: 'Get all schedules',
      parameters: ['status', 'limit', 'offset'],
    },
    {
      method: 'GET',
      path: '/api/v1/schedules/:id',
      description: 'Get specific schedule by ID',
      parameters: ['id'],
    },
    {
      method: 'POST',
      path: '/api/v1/schedules',
      description: 'Create new schedule',
      parameters: ['clientName', 'shiftTime', 'location'],
    },
    {
      method: 'GET',
      path: '/api/v1/schedules/stats',
      description: 'Get schedule statistics',
      parameters: [],
    },
    {
      method: 'GET',
      path: '/docs/swagger.json',
      description: 'OpenAPI specification',
      parameters: [],
    },
  ];

  const allEndpoints = endpoints.length > 0 ? endpoints : defaultEndpoints;

  const testEndpoint = async (endpoint: APIEndpoint) => {
    setIsLoading(true);
    setResponse(null);

    try {
      const startTime = performance.now();
      const url = `http://localhost:9090${endpoint.path}`;

      const response = await fetch(url, {
        method: endpoint.method,
        headers: {
          'Content-Type': 'application/json',
        },
      });

      const endTime = performance.now();
      const responseTime = Math.round(endTime - startTime);

      let data;
      try {
        data = await response.json();
      } catch (error) {
        data = null;
      }

      setResponse({
        status: response.status,
        data: response.ok ? data : undefined,
        error: response.ok ? undefined : `HTTP ${response.status}: ${response.statusText}`,
        responseTime,
        timestamp: new Date().toISOString(),
      });
    } catch (error) {
      const endTime = performance.now();
      setResponse({
        status: 0,
        error: error instanceof Error ? error.message : 'Unknown error',
        responseTime: Math.round(performance.now() - endTime),
        timestamp: new Date().toISOString(),
      });
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
  };

  const copyCURL = (endpoint: APIEndpoint) => {
    const curlCommand = `curl -X ${endpoint.method} http://localhost:9090${endpoint.path} -H "Content-Type: application/json"`;
    copyToClipboard(curlCommand);
  };

  return (
    <div className="bg-white rounded-lg shadow-xl border border-gray-200">
      <div className="p-6">
        <h3 className="text-xl font-semibold text-gray-900 mb-4">API Endpoint Tester</h3>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Endpoint List */}
          <div className="lg:col-span-1">
            <h4 className="font-medium text-gray-900 mb-3">Available Endpoints</h4>
            <div className="space-y-2">
              {allEndpoints.map((endpoint, index) => (
                <button
                  key={index}
                  onClick={() => setSelectedEndpoint(endpoint)}
                  className={`w-full text-left p-3 rounded-lg border transition-colors ${
                    selectedEndpoint === endpoint
                      ? 'border-blue-300 bg-blue-50'
                      : 'border-gray-200 hover:bg-gray-50'
                  }`}
                >
                  <div className="flex items-center justify-between">
                    <span className="font-mono text-sm font-semibold text-gray-900">
                      {endpoint.method}
                    </span>
                    <span className="font-mono text-xs text-gray-600 bg-gray-100 px-2 py-1 rounded">
                      {endpoint.path}
                    </span>
                  </div>
                  <p className="text-xs text-gray-600 mt-1">{endpoint.description}</p>
                </button>
              ))}
            </div>
          </div>

          {/* Endpoint Details & Tester */}
          <div className="lg:col-span-2">
            {selectedEndpoint ? (
              <div className="space-y-4">
                {/* Endpoint Header */}
                <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                  <div className="flex items-center space-x-3">
                    <span className="font-mono text-lg font-bold text-gray-900">
                      {selectedEndpoint.method}
                    </span>
                    <span className="font-mono text-gray-800">{selectedEndpoint.path}</span>
                  </div>
                  <button
                    onClick={() => copyCURL(selectedEndpoint)}
                    className="px-3 py-1 text-sm bg-gray-100 hover:bg-gray-200 rounded-md flex items-center space-x-1"
                  >
                    <Copy className="h-3 w-3" />
                    <span>CURL</span>
                  </button>
                </div>

                {/* Test Button */}
                <div className="flex items-center space-x-3">
                  <button
                    onClick={() => testEndpoint(selectedEndpoint)}
                    disabled={isLoading}
                    className="flex items-center space-x-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md disabled:opacity-50"
                  >
                    {isLoading ? (
                      <RefreshCw className="h-4 w-4 animate-spin" />
                    ) : (
                      <Play className="h-4 w-4" />
                    )}
                    <span>Test Endpoint</span>
                  </button>

                  {selectedEndpoint.example && (
                    <button
                      onClick={() => setCustomParams(JSON.stringify(selectedEndpoint.example, null, 2))}
                      className="px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 rounded-md"
                    >
                      Use Example
                    </button>
                  )}
                </div>

                {/* Response */}
                {response && (
                  <div className={`rounded-lg border p-4 ${
                    response.status >= 400
                      ? 'border-red-200 bg-red-50'
                      : 'border-green-200 bg-green-50'
                  }`}>
                    <div className="flex items-center justify-between mb-3">
                      <div className="flex items-center space-x-2">
                        {response.status >= 400 ? (
                          <AlertCircle className="h-5 w-5 text-red-500" />
                        ) : (
                          <CheckCircle className="h-5 w-5 text-green-500" />
                        )}
                        <span className="font-medium">
                          {response.status} {response.status >= 400 ? 'Error' : 'Success'}
                        </span>
                        <span className="text-sm text-gray-600">
                          {response.responseTime}ms
                        </span>
                      </div>
                      <span className="text-xs text-gray-500">
                        {new Date(response.timestamp).toLocaleTimeString()}
                      </span>
                    </div>

                    {response.error && (
                      <div className="mb-3">
                        <h4 className="text-sm font-medium text-red-800 mb-1">Error:</h4>
                        <p className="text-sm text-red-700">{response.error}</p>
                      </div>
                    )}

                    {response.data && (
                      <div>
                        <h4 className="text-sm font-medium text-gray-800 mb-2">Response:</h4>
                        <div className="bg-black/5 p-3 rounded overflow-x-auto">
                          <pre className="text-xs">
                            <code>{JSON.stringify(response.data, null, 2)}</code>
                          </pre>
                        </div>
                      </div>
                    )}

                    <div className="mt-3 flex space-x-2">
                      <button
                        onClick={() => copyToClipboard(JSON.stringify(response.data || response.error, null, 2))}
                        className="text-xs px-3 py-1 bg-gray-100 hover:bg-gray-200 rounded flex items-center space-x-1"
                      >
                        <Copy className="h-3 w-3" />
                        <span>Copy Response</span>
                      </button>
                    </div>
                  </div>
                )}

                {/* API Documentation */}
                <div className="mt-6">
                  <h4 className="font-medium text-gray-900 mb-2">Documentation</h4>
                  <p className="text-sm text-gray-600">{selectedEndpoint.description}</p>

                  {selectedEndpoint.parameters && selectedEndpoint.parameters.length > 0 && (
                    <div className="mt-3">
                      <h5 className="text-sm font-medium text-gray-900 mb-2">Parameters:</h5>
                      <ul className="text-sm text-gray-600 space-y-1">
                        {selectedEndpoint.parameters.map((param, index) => (
                          <li key={index} className="font-mono bg-gray-100 px-2 py-1 rounded">
                            {param}
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              </div>
            ) : (
              <div className="text-center py-12 text-gray-500">
                <AlertCircle className="h-12 w-12 mx-auto mb-3 opacity-50" />
                <p>Select an endpoint from the list to test it</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}