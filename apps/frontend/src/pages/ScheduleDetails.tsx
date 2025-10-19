import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  ArrowLeft,
  Clock,
  MapPin,
  Play,
  Square,
  CheckCircle,
  XCircle,
  AlertTriangle,
  Activity,
  Loader2
} from "lucide-react";
import { mockApi } from "../api/standalone";
import { type Schedule, type Location } from "../api/standalone";

export function ScheduleDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const [visitStartTime, setVisitStartTime] = useState<Date | null>(null);
  const [visitEndTime, setVisitEndTime] = useState<Date | null>(null);
  const [startLocation, setStartLocation] = useState<Location | null>(null);
  const [endLocation, setEndLocation] = useState<Location | null>(null);
  const [isGeolocationError, setIsGeolocationError] = useState(false);

  // Fetch schedule details
  const { data: schedule, isLoading } = useQuery({
    queryKey: ["schedule", id],
    queryFn: () => id ? mockApi.getScheduleById(id) : Promise.reject("No schedule ID"),
    enabled: !!id,
  });

  // Start visit mutation
  const startVisitMutation = useMutation({
    mutationFn: async () => {
      if (!id || !startLocation) throw new Error("Missing schedule ID or location");

      await mockApi.startVisit(id, startLocation);
      setVisitStartTime(new Date());
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["schedule", id] });
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
      queryClient.invalidateQueries({ queryKey: ["stats"] });
    },
    onError: (error) => {
      console.error("Failed to start visit:", error);
    },
  });

  // End visit mutation
  const endVisitMutation = useMutation({
    mutationFn: async () => {
      if (!id || !endLocation) throw new Error("Missing schedule ID or location");

      await mockApi.endVisit(id, endLocation);
      setVisitEndTime(new Date());
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["schedule", id] });
      queryClient.invalidateQueries({ queryKey: ["schedules"] });
      queryClient.invalidateQueries({ queryKey: ["stats"] });
    },
    onError: (error) => {
      console.error("Failed to end visit:", error);
    },
  });

  // Update task mutation
  const updateTaskMutation = useMutation({
    mutationFn: async ({ taskId, status, reason }: { taskId: string; status: string; reason?: string }) => {
      await mockApi.updateTask(taskId, { status, reason });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["schedule", id] });
    },
  });

  // Get current geolocation
  const getCurrentLocation = (): Promise<Location> => {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error("Geolocation is not supported by this browser."));
        return;
      }

      navigator.geolocation.getCurrentPosition(
        (position) => {
          const location: Location = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            accuracy: position.coords.accuracy,
            timestamp: new Date(),
          };
          resolve(location);
        },
        (error) => {
          reject(error);
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 0,
        }
      );
    });
  };

  const handleStartVisit = async () => {
    try {
      setIsGeolocationError(false);
      const location = await getCurrentLocation();
      setStartLocation(location);
      startVisitMutation.mutate();
    } catch (error) {
      console.error("Geolocation error:", error);
      setIsGeolocationError(true);
    }
  };

  const handleEndVisit = async () => {
    try {
      setIsGeolocationError(false);
      const location = await getCurrentLocation();
      setEndLocation(location);
      endVisitMutation.mutate();
    } catch (error) {
      console.error("Geolocation error:", error);
      setIsGeolocationError(true);
    }
  };

  const handleTaskUpdate = (taskId: string, status: string, reason?: string) => {
    updateTaskMutation.mutate({ taskId, status, reason });
  };

  if (isLoading) {
    return (
      <div className="flex items-center justify-center h-64">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
        <span className="ml-2 text-gray-600">Loading schedule details...</span>
      </div>
    );
  }

  if (!schedule) {
    return (
      <div className="text-center py-12">
        <AlertTriangle className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">Schedule not found</h3>
        <p className="mt-1 text-sm text-gray-500">The requested schedule does not exist.</p>
        <div className="mt-4">
          <Link
            to="/"
            className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
          >
            Back to Schedules
          </Link>
        </div>
      </div>
    );
  }

  const isVisitStarted = schedule.visit?.status === "in_progress" || schedule.visit?.status === "completed";
  const isVisitCompleted = schedule.visit?.status === "completed";

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center space-x-4">
        <Link
          to="/"
          className="flex items-center text-gray-600 hover:text-gray-900"
        >
          <ArrowLeft className="h-5 w-5 mr-1" />
          Back to Schedules
        </Link>
        <div className="flex items-center space-x-2">
          <h1 className="text-2xl font-bold text-gray-900">{schedule.clientName}</h1>
          <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
            {schedule.shiftTime}
          </span>
        </div>
      </div>

      {/* Visit Status */}
      <div className="bg-white rounded-lg shadow p-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-medium text-gray-900">Visit Status</h2>
          {isVisitCompleted && schedule.visit && (
            <div className="text-sm text-gray-600">
              Duration: {Math.round((new Date(schedule.visit.endTime).getTime() - new Date(schedule.visit.startTime).getTime()) / 60000)} minutes
            </div>
          )}
        </div>

        <div className="space-y-4">
          {isGeolocationError && (
            <div className="bg-red-50 border border-red-200 rounded-md p-4">
              <div className="flex">
                <XCircle className="h-5 w-5 text-red-400" />
                <div className="ml-3">
                  <h3 className="text-sm font-medium text-red-800">Geolocation Error</h3>
                  <p className="mt-1 text-sm text-red-700">
                    Unable to get your location. Please check your browser permissions and try again.
                  </p>
                </div>
              </div>
            </div>
          )}

          {!isVisitStarted ? (
            <button
              onClick={handleStartVisit}
              disabled={startVisitMutation.isPending}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
            >
              {startVisitMutation.isPending ? (
                <>
                  <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4" />
                  Starting...
                </>
              ) : (
                <>
                  <Play className="-ml-1 mr-2 h-4 w-4" />
                  Start Visit
                </>
              )}
            </button>
          ) : isVisitCompleted ? (
            <div className="flex items-center space-x-2">
              <CheckCircle className="h-5 w-5 text-green-500" />
              <span className="text-green-700 font-medium">Visit Completed</span>
            </div>
          ) : (
            <button
              onClick={handleEndVisit}
              disabled={endVisitMutation.isPending}
              className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
            >
              {endVisitMutation.isPending ? (
                <>
                  <Loader2 className="animate-spin -ml-1 mr-2 h-4 w-4" />
                  Ending...
                </>
              ) : (
                <>
                  <Square className="-ml-1 mr-2 h-4 w-4" />
                  End Visit
                </>
              )}
            </button>
          )}

          {schedule.visit && (
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-6">
              <div className="border rounded-lg p-4">
                <h4 className="text-sm font-medium text-gray-900 mb-2">Start Location</h4>
                <div className="space-y-1 text-sm text-gray-600">
                  <div className="flex items-center">
                    <Clock className="h-4 w-4 mr-2" />
                    {new Date(schedule.visit.startTime).toLocaleString()}
                  </div>
                  <div className="flex items-center">
                    <MapPin className="h-4 w-4 mr-2" />
                    {schedule.visit.startLocation.latitude.toFixed(6)}, {schedule.visit.startLocation.longitude.toFixed(6)}
                  </div>
                  {schedule.visit.startLocation.accuracy && (
                    <div className="text-xs">
                      Accuracy: ±{Math.round(schedule.visit.startLocation.accuracy)}m
                    </div>
                  )}
                </div>
              </div>

              {schedule.visit.endTime && (
                <div className="border rounded-lg p-4">
                  <h4 className="text-sm font-medium text-gray-900 mb-2">End Location</h4>
                  <div className="space-y-1 text-sm text-gray-600">
                    <div className="flex items-center">
                      <Clock className="h-4 w-4 mr-2" />
                      {new Date(schedule.visit.endTime).toLocaleString()}
                    </div>
                    <div className="flex items-center">
                      <MapPin className="h-4 w-4 mr-2" />
                      {schedule.visit.endLocation!.latitude.toFixed(6)}, {schedule.visit.endLocation!.longitude.toFixed(6)}
                    </div>
                    {schedule.visit.endLocation?.accuracy && (
                      <div className="text-xs">
                        Accuracy: ±{Math.round(schedule.visit.endLocation.accuracy)}m
                      </div>
                    )}
                  </div>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Tasks */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <h2 className="text-lg font-medium text-gray-900">Care Activities</h2>
          <p className="mt-1 text-sm text-gray-500">Mark tasks as completed during your visit</p>
        </div>
        <div className="divide-y divide-gray-200">
          {schedule.tasks.length === 0 ? (
            <div className="px-6 py-12 text-center">
              <Activity className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-sm font-medium text-gray-900">No tasks assigned</h3>
              <p className="mt-1 text-sm text-gray-500">No care activities have been assigned for this schedule.</p>
            </div>
          ) : (
            schedule.tasks.map((task) => (
              <div key={task.id} className="px-6 py-4">
                <div className="flex items-center justify-between">
                  <div className="flex-1">
                    <h4 className="text-sm font-medium text-gray-900">{task.name}</h4>
                    {task.description && (
                      <p className="mt-1 text-sm text-gray-500">{task.description}</p>
                    )}
                    {task.reason && (
                      <p className="mt-1 text-sm text-red-600">Reason: {task.reason}</p>
                    )}
                  </div>
                  <div className="flex items-center space-x-2">
                    {task.status === "completed" ? (
                      <CheckCircle className="h-5 w-5 text-green-500" />
                    ) : task.status === "not_completed" ? (
                      <XCircle className="h-5 w-5 text-red-500" />
                    ) : (
                      <div className="flex space-x-2">
                        <button
                          onClick={() => handleTaskUpdate(task.id, "completed")}
                          disabled={updateTaskMutation.isPending}
                          className="inline-flex items-center px-3 py-1 border border-transparent text-xs font-medium rounded text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
                        >
                          Complete
                        </button>
                        <button
                          onClick={() => {
                            const reason = prompt("Please provide a reason for not completing this task:");
                            if (reason) {
                              handleTaskUpdate(task.id, "not_completed", reason);
                            }
                          }}
                          disabled={updateTaskMutation.isPending}
                          className="inline-flex items-center px-3 py-1 border border-gray-300 text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50"
                        >
                          Not Completed
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}