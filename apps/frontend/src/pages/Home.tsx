import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import {
  Calendar,
  Clock,
  MapPin,
  CheckCircle,
  XCircle,
  Activity,
  AlertCircle,
  Monitor,
  Code
} from "lucide-react";
// StatsCard component moved inline to avoid dependency issues
function StatsCard({ title, value, description, icon: Icon, className = "" }) {
  return (
    <div className={`${className} bg-white rounded-lg shadow p-6`}>
      <div className="flex items-center justify-between pb-2">
        <h3 className="text-sm font-medium">{title}</h3>
        <Icon className="h-4 w-4 text-gray-500" />
      </div>
      <div>
        <div className="text-2xl font-bold">{value}</div>
        {description && <p className="text-xs text-gray-500">{description}</p>}
      </div>
    </div>
  );
}
import { mockApi } from "../api/standalone";
import { type Schedule, type Stats } from "../api/standalone";
import { APIResponseViewer } from "../components/APIResponseViewer";
import { APIDocsTester } from "../components/APIDocsTester";

export function Home() {
  const [selectedStatus, setSelectedStatus] = useState<string>("all");
  const [showAPIViewer, setShowAPIViewer] = useState(false);
  const [showAPIDocs, setShowAPIDocs] = useState(false);

  // Fetch stats
  const { data: stats, isLoading: statsLoading } = useQuery({
    queryKey: ["stats"],
    queryFn: () => mockApi.getStats(),
  });

  // Fetch schedules
  const { data: schedules, isLoading: schedulesLoading } = useQuery({
    queryKey: ["schedules"],
    queryFn: () => mockApi.getSchedules(),
  });

  // Filter schedules based on selected status
  const filteredSchedules = schedules?.filter(schedule => {
    if (selectedStatus === "all") return true;
    return schedule.status === selectedStatus;
  });

  // Status badge component
  const StatusBadge = ({ status }: { status: Schedule["status"] }) => {
    const variants = {
      upcoming: "bg-blue-100 text-blue-800",
      in_progress: "bg-yellow-100 text-yellow-800",
      completed: "bg-green-100 text-green-800",
      missed: "bg-red-100 text-red-800",
    };

    const icons = {
      upcoming: Clock,
      in_progress: Activity,
      completed: CheckCircle,
      missed: XCircle,
    };

    const Icon = icons[status];
    const label = status.replace("_", " ");

    return (
      <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${variants[status]}`}>
        <Icon className="w-3 h-3 mr-1" />
        {label}
      </span>
    );
  };

  return (
    <div className="space-y-8">
      {/* Stats Dashboard */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatsCard
          title="Total Schedules"
          value={stats?.total || 0}
          description="All scheduled visits"
          icon={Calendar}
          className="col-span-1"
        />
        <StatsCard
          title="Upcoming"
          value={stats?.upcoming || 0}
          description="For today"
          icon={Clock}
          className="col-span-1"
        />
        <StatsCard
          title="Completed"
          value={stats?.completed || 0}
          description="Successfully finished"
          icon={CheckCircle}
          className="col-span-1"
        />
        <StatsCard
          title="Missed"
          value={stats?.missed || 0}
          description="Not completed"
          icon={XCircle}
          className="col-span-1"
        />
      </div>

      {/* API Monitoring Controls */}
      <div className="bg-white rounded-lg shadow p-6 border border-gray-200">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-lg font-medium text-gray-900 flex items-center space-x-2">
            <Monitor className="h-5 w-5 text-blue-600" />
            <span>API Monitoring Tools</span>
          </h2>
          <div className="flex space-x-2">
            <button
              onClick={() => setShowAPIViewer(!showAPIViewer)}
              className="flex items-center space-x-1 px-3 py-2 text-sm bg-blue-100 hover:bg-blue-200 text-blue-700 rounded-md"
            >
              <Monitor className="h-3 w-3" />
              <span>{showAPIViewer ? 'Hide' : 'Show'} API Viewer</span>
            </button>
            <button
              onClick={() => setShowAPIDocs(!showAPIDocs)}
              className="flex items-center space-x-1 px-3 py-2 text-sm bg-purple-100 hover:bg-purple-200 text-purple-700 rounded-md"
            >
              <Code className="h-3 w-3" />
              <span>{showAPIDocs ? 'Hide' : 'Show'} API Tester</span>
            </button>
          </div>
        </div>

        {showAPIViewer && <APIResponseViewer showOnPage={true} />}
        {showAPIDocs && <APIDocsTester />}
      </div>

      {/* Schedule List */}
      <div className="bg-white rounded-lg shadow">
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <h2 className="text-lg font-medium text-gray-900">Today's Schedules</h2>
            <div className="mt-2 sm:mt-0">
              <select
                value={selectedStatus}
                onChange={(e) => setSelectedStatus(e.target.value)}
                className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
              >
                <option value="all">All Statuses</option>
                <option value="upcoming">Upcoming</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
                <option value="missed">Missed</option>
              </select>
            </div>
          </div>
        </div>

        <div className="divide-y divide-gray-200">
          {schedulesLoading ? (
            <div className="px-6 py-12">
              <div className="space-y-4">
                {[1, 2, 3].map((i) => (
                  <div key={i} className="animate-pulse">
                    <div className="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                    <div className="h-4 bg-gray-200 rounded w-1/2"></div>
                  </div>
                ))}
              </div>
            </div>
          ) : filteredSchedules && filteredSchedules.length > 0 ? (
            filteredSchedules.map((schedule) => (
              <div key={schedule.id} className="px-6 py-4 hover:bg-gray-50">
                <div className="flex items-center justify-between">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center space-x-3">
                      <h3 className="text-sm font-medium text-gray-900 truncate">
                        {schedule.clientName}
                      </h3>
                      <StatusBadge status={schedule.status} />
                    </div>
                    <div className="mt-1 flex items-center text-sm text-gray-500">
                      <Clock className="w-4 h-4 mr-1" />
                      <span>{schedule.shiftTime}</span>
                      <MapPin className="w-4 h-4 ml-3 mr-1" />
                      <span>{schedule.location}</span>
                    </div>
                  </div>
                  <div className="ml-4 flex-shrink-0">
                    <Link
                      to={`/schedules/${schedule.id}`}
                      className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                    >
                      View Details
                    </Link>
                  </div>
                </div>
              </div>
            ))
          ) : (
            <div className="px-6 py-12 text-center">
              <AlertCircle className="mx-auto h-12 w-12 text-gray-400" />
              <h3 className="mt-2 text-sm font-medium text-gray-900">No schedules found</h3>
              <p className="mt-1 text-sm text-gray-500">
                {selectedStatus === "all"
                  ? "There are no schedules for today."
                  : `There are no schedules with "${selectedStatus}" status.`
                }
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}