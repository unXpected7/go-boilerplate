import { apiContract } from "@boilerplate/openapi/contracts";
import type { ServerInferRequest } from "@ts-rest/core";

export type TRequests = ServerInferRequest<typeof apiContract>;

// EVV Specific Types
export interface Location {
  latitude: number;
  longitude: number;
  accuracy?: number;
  timestamp?: Date;
}

export interface Visit {
  id: string;
  scheduleId: string;
  startTime: Date;
  endTime?: Date;
  startLocation: Location;
  endLocation?: Location;
  status: "not_started" | "in_progress" | "completed";
  createdAt: Date;
  updatedAt: Date;
}

export interface Task {
  id: string;
  scheduleId: string;
  name: string;
  description?: string;
  status: "pending" | "completed" | "not_completed";
  reason?: string;
  completedAt?: Date;
  createdAt: Date;
  updatedAt: Date;
}

export interface Schedule {
  id: string;
  clientName: string;
  shiftTime: string;
  location: string;
  status: "missed" | "upcoming" | "in_progress" | "completed";
  visit?: Visit;
  tasks: Task[];
  createdAt: Date;
  updatedAt: Date;
}

export interface Stats {
  total: number;
  missed: number;
  upcoming: number;
  completed: number;
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// API Request Types
export interface StartVisitRequest {
  scheduleId: string;
  location: Location;
}

export interface EndVisitRequest {
  scheduleId: string;
  location: Location;
}

export interface UpdateTaskRequest {
  taskId: string;
  status: "completed" | "not_completed";
  reason?: string;
}
