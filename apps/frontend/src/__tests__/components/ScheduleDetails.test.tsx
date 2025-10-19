import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { ScheduleDetails } from '../../pages/ScheduleDetails';
import { mockApi } from '../../api/standalone';

// Mock the mockApi module
jest.mock('../../api/standalone');
const mockedMockApi = mockApi as jest.Mocked<typeof mockApi>;

// Mock useParams hook
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  useParams: () => ({ id: '1' }),
}));

describe('ScheduleDetails Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render loading state initially', () => {
    mockedMockApi.getScheduleById.mockReturnValue(new Promise(() => {}));

    render(<ScheduleDetails />);

    expect(screen.getByText('Loading schedule details...')).toBeInTheDocument();
  });

  it('should render schedule details when data loads', async () => {
    const mockSchedule = {
      id: '1',
      clientName: 'John Smith',
      shiftTime: '09:00 - 12:00',
      location: '123 Main St',
      status: 'upcoming',
      tasks: [
        {
          id: 'task1',
          name: 'Morning Medication',
          description: 'Administer morning medication',
          status: 'pending',
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    mockedMockApi.getScheduleById.mockResolvedValue(mockSchedule);

    render(<ScheduleDetails />);

    await waitFor(() => {
      expect(screen.getByText('John Smith')).toBeInTheDocument();
      expect(screen.getByText('09:00 - 12:00')).toBeInTheDocument();
      expect(screen.getByText('Morning Medication')).toBeInTheDocument();
    });
  });

  it('should show "Schedule not found" when schedule does not exist', async () => {
    mockedMockApi.getScheduleById.mockRejectedValue(new Error('Schedule not found'));

    render(<ScheduleDetails />);

    await waitFor(() => {
      expect(screen.getByText('Schedule not found')).toBeInTheDocument();
    });
  });

  it('should render start visit button when visit not started', async () => {
    const mockSchedule = {
      id: '1',
      clientName: 'John Smith',
      shiftTime: '09:00 - 12:00',
      location: '123 Main St',
      status: 'upcoming',
      tasks: [],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    mockedMockApi.getScheduleById.mockResolvedValue(mockSchedule);

    render(<ScheduleDetails />);

    await waitFor(() => {
      expect(screen.getByText('Start Visit')).toBeInTheDocument();
    });
  });

  it('should render task management when visit is in progress', async () => {
    const mockSchedule = {
      id: '1',
      clientName: 'John Smith',
      shiftTime: '09:00 - 12:00',
      location: '123 Main St',
      status: 'in_progress',
      visit: {
        id: 'visit1',
        scheduleId: '1',
        startTime: new Date(),
        status: 'in_progress',
        startLocation: {
          latitude: 40.7128,
          longitude: -74.0060,
          timestamp: new Date(),
        },
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      tasks: [
        {
          id: 'task1',
          name: 'Morning Medication',
          description: 'Administer morning medication',
          status: 'pending',
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    mockedMockApi.getScheduleById.mockResolvedValue(mockSchedule);

    render(<ScheduleDetails />);

    await waitFor(() => {
      expect(screen.getByText('Morning Medication')).toBeInTheDocument();
      expect(screen.getByText('Complete')).toBeInTheDocument();
      expect(screen.getByText('Not Completed')).toBeInTheDocument();
    });
  });

  it('should handle task completion', async () => {
    const mockSchedule = {
      id: '1',
      clientName: 'John Smith',
      shiftTime: '09:00 - 12:00',
      location: '123 Main St',
      status: 'in_progress',
      visit: {
        id: 'visit1',
        scheduleId: '1',
        startTime: new Date(),
        status: 'in_progress',
        startLocation: {
          latitude: 40.7128,
          longitude: -74.0060,
          timestamp: new Date(),
        },
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      tasks: [
        {
          id: 'task1',
          name: 'Morning Medication',
          description: 'Administer morning medication',
          status: 'pending',
          createdAt: new Date(),
          updatedAt: new Date(),
        },
      ],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    mockedMockApi.getScheduleById.mockResolvedValue(mockSchedule);
    mockedMockApi.updateTask.mockResolvedValue(undefined);

    render(<ScheduleDetails />);

    await waitFor(() => {
      expect(screen.getByText('Morning Medication')).toBeInTheDocument();
    });

    const completeButton = screen.getByText('Complete');
    fireEvent.click(completeButton);

    await waitFor(() => {
      expect(mockedMockApi.updateTask).toHaveBeenCalledWith('task1', {
        status: 'completed',
      });
    });
  });

  it('should handle geolocation error', async () => {
    const mockSchedule = {
      id: '1',
      clientName: 'John Smith',
      shiftTime: '09:00 - 12:00',
      location: '123 Main St',
      status: 'upcoming',
      tasks: [],
      createdAt: new Date(),
      updatedAt: new Date(),
    };

    mockedMockApi.getScheduleById.mockResolvedValue(mockSchedule);

    // Mock geolocation error
    const originalGeolocation = navigator.geolocation;
    Object.defineProperty(navigator, 'geolocation', {
      value: undefined,
      configurable: true,
    });

    render(<ScheduleDetails />);

    await waitFor(() => {
      expect(screen.getByText('Geolocation Error')).toBeInTheDocument();
    });

    // Restore original geolocation
    Object.defineProperty(navigator, 'geolocation', {
      value: originalGeolocation,
      configurable: true,
    });
  });
});