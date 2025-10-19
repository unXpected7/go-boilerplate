import { API_URL } from "@/config/env";
import { apiContract } from "@boilerplate/openapi/contracts";
import { useAuth } from "@clerk/clerk-react";
import { initClient } from "@ts-rest/core";
import axios, {
  type Method,
  AxiosError,
  isAxiosError,
  type AxiosResponse,
} from "axios";
import {
  type Schedule,
  type Stats,
  type StartVisitRequest,
  type EndVisitRequest,
  type UpdateTaskRequest,
  type Location,
  type ApiResponse,
  type PaginatedResponse,
} from "./types";

type Headers = Awaited<
  ReturnType<NonNullable<Parameters<typeof initClient>[1]["api"]>>
>["headers"];

export type TApiClient = ReturnType<typeof useApiClient>;

export const useApiClient = ({ isBlob = false }: { isBlob?: boolean } = {}) => {
  const { getToken } = useAuth();

  return initClient(apiContract, {
    baseUrl: "",
    baseHeaders: {
      "Content-Type": "application/json",
    },
    api: async ({ path, method, headers, body }) => {
      const token = await getToken({ template: "custom" });

      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const makeRequest = async (retryCount = 0): Promise<any> => {
        try {
          const result = await axios.request({
            method: method as Method,
            url: `${API_URL}/api${path}`,
            headers: {
              ...headers,
              ...(token ? { Authorization: `Bearer ${token}` } : {}),
            },
            data: body,
            ...(isBlob ? { responseType: "blob" } : {}),
          });
          return {
            status: result.status,
            body: result.data,
            headers: result.headers as unknown as Headers,
          };
          // eslint-disable-next-line @typescript-eslint/no-explicit-any
        } catch (e: Error | AxiosError | any) {
          if (isAxiosError(e)) {
            const error = e as AxiosError;
            const response = error.response as AxiosResponse;

            // If unauthorized and we haven't retried yet, retry
            if (response?.status === 401 && retryCount < 2) {
              return makeRequest(retryCount + 1);
            }

            return {
              status: response?.status || 500,
              body: response?.data || { message: "Internal server error" },
              headers: (response?.headers as unknown as Headers) || {},
            };
          }
          throw e;
        }
      };

      return makeRequest();
    },
  });
};

// EVV Specific API Functions
export const evvApi = {
  // Schedule APIs
  getSchedules: async (client: TApiClient): Promise<Schedule[]> => {
    const result = await client.schedules.get({});
    return result.body.data;
  },

  getSchedulesToday: async (client: TApiClient): Promise<Schedule[]> => {
    const result = await client.schedules.getToday({});
    return result.body.data;
  },

  getScheduleById: async (client: TApiClient, id: string): Promise<Schedule> => {
    const result = await client.schedules.getById({ path: { id } });
    return result.body.data;
  },

  // Visit APIs
  startVisit: async (
    client: TApiClient,
    scheduleId: string,
    location: Location
  ): Promise<void> => {
    await client.schedules.startVisit({
      path: { id: scheduleId },
      body: { location },
    });
  },

  endVisit: async (
    client: TApiClient,
    scheduleId: string,
    location: Location
  ): Promise<void> => {
    await client.schedules.endVisit({
      path: { id: scheduleId },
      body: { location },
    });
  },

  // Task APIs
  updateTask: async (
    client: TApiClient,
    taskId: string,
    data: UpdateTaskRequest
  ): Promise<void> => {
    await client.tasks.updateTask({
      path: { taskId },
      body: data,
    });
  },

  // Stats API
  getStats: async (client: TApiClient): Promise<Stats> => {
    const result = await client.schedules.getStats({});
    return result.body.data;
  },
};
