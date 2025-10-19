// Standalone API client for testing without workspace dependencies
// import axios, { type AxiosResponse } from "axios"; // Commented out since not used in mock

const API_URL = "http://localhost:8080/api";

// Types moved here to avoid dependency issues
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

// Mock API functions for testing
export const mockApi = {
  // Mock data for development
  getMockSchedules: (): Schedule[] => [
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
        {
          id: "task2",
          scheduleId: "1",
          name: "Assist with Bathing",
          description: "Help with morning bathing routine",
          status: "pending",
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    {
      id: "2",
      clientName: "Mary Johnson",
      shiftTime: "14:00 - 17:00",
      location: "456 Oak Ave, Somewhere",
      status: "upcoming",
      tasks: [
        {
          id: "task3",
          scheduleId: "2",
          name: "Afternoon Meal Preparation",
          description: "Prepare and serve afternoon meal",
          status: "pending",
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    },
    {
      id: "3",
      clientName: "Robert Brown",
      shiftTime: "10:00 - 11:00",
      location: "789 Pine St, Nowhere",
      status: "completed",
      visit: {
        id: "visit1",
        scheduleId: "3",
        startTime: new Date(Date.now() - 60000), // 1 minute ago
        endTime: new Date(),
        startLocation: {
          latitude: 40.7128,
          longitude: -74.0060,
          accuracy: 10,
          timestamp: new Date(Date.now() - 60000),
        },
        endLocation: {
          latitude: 40.7129,
          longitude: -74.0061,
          accuracy: 15,
          timestamp: new Date(),
        },
        status: "completed",
        createdAt: new Date(Date.now() - 60000),
        updatedAt: new Date(),
      },
      tasks: [
        {
          id: "task4",
          scheduleId: "3",
          name: "Physical Therapy",
          description: "Assist with physical therapy exercises",
          status: "completed",
          completedAt: new Date(),
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    },
  ],

  getMockStats: (): Stats => ({
    total: 8,
    missed: 2,
    upcoming: 3,
    completed: 3,
  }),

  // API functions
  getSchedules: async (): Promise<Schedule[]> => {
    // Simulate API delay
    await new Promise(resolve => setTimeout(resolve, 500));
    return mockApi.getMockSchedules();
  },

  getSchedulesToday: async (): Promise<Schedule[]> => {
    await new Promise(resolve => setTimeout(resolve, 300));
    return mockApi.getMockSchedules().filter(s => s.status !== "missed");
  },

  getScheduleById: async (id: string): Promise<Schedule> => {
    await new Promise(resolve => setTimeout(resolve, 300));
    const schedule = mockApi.getMockSchedules().find(s => s.id === id);
    if (!schedule) {
      throw new Error("Schedule not found");
    }
    return schedule;
  },

  startVisit: async (scheduleId: string, location: Location): Promise<void> => {
    await new Promise(resolve => setTimeout(resolve, 1000));
    console.log(`Starting visit for schedule ${scheduleId} at`, location);
  },

  endVisit: async (scheduleId: string, location: Location): Promise<void> => {
    await new Promise(resolve => setTimeout(resolve, 1000));
    console.log(`Ending visit for schedule ${scheduleId} at`, location);
  },

  updateTask: async (taskId: string, data: { status: string; reason?: string }): Promise<void> => {
    await new Promise(resolve => setTimeout(resolve, 500));
    console.log(`Updating task ${taskId} to ${data.status}`, data.reason);
  },

  getStats: async (): Promise<Stats> => {
    await new Promise(resolve => setTimeout(resolve, 300));
    return mockApi.getMockStats();
  },
};