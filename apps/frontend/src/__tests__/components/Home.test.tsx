import { render, screen, waitFor, fireEvent } from '@testing-library/react';
import { Home } from '../../pages/Home';
import { mockApi } from '../../api/standalone';

// Mock the mockApi module
jest.mock('../../api/standalone');
const mockedMockApi = mockApi as jest.Mocked<typeof mockApi>;

describe('Home Component', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render loading state initially', () => {
    mockedMockApi.getSchedules.mockReturnValue(new Promise(() => {}));
    mockedMockApi.getStats.mockReturnValue(new Promise(() => {}));

    render(<Home />);

    expect(screen.getByText('Loading...')).toBeInTheDocument();
  });

  it('should render stats cards when data loads', async () => {
    mockedMockApi.getStats.mockResolvedValue({
      total: 8,
      missed: 2,
      upcoming: 3,
      completed: 3,
    });

    mockedMockApi.getSchedules.mockResolvedValue([
      {
        id: '1',
        clientName: 'John Smith',
        shiftTime: '09:00 - 12:00',
        location: '123 Main St',
        status: 'upcoming',
        tasks: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    ]);

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText('8')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
      expect(screen.getByText('3')).toBeInTheDocument();
      expect(screen.getByText('2')).toBeInTheDocument();
    });
  });

  it('should render schedule list with status badges', async () => {
    mockedMockApi.getSchedules.mockResolvedValue([
      {
        id: '1',
        clientName: 'John Smith',
        shiftTime: '09:00 - 12:00',
        location: '123 Main St',
        status: 'upcoming',
        tasks: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '2',
        clientName: 'Mary Johnson',
        shiftTime: '14:00 - 17:00',
        location: '456 Oak Ave',
        status: 'completed',
        tasks: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    ]);

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText('John Smith')).toBeInTheDocument();
      expect(screen.getByText('Mary Johnson')).toBeInTheDocument();
    });

    const upcomingBadges = screen.getAllByText('Upcoming');
    const completedBadges = screen.getAllByText('Completed');

    expect(upcomingBadges).toHaveLength(1);
    expect(completedBadges).toHaveLength(1);
  });

  it('should filter schedules by status', async () => {
    mockedMockApi.getSchedules.mockResolvedValue([
      {
        id: '1',
        clientName: 'John Smith',
        shiftTime: '09:00 - 12:00',
        location: '123 Main St',
        status: 'upcoming',
        tasks: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      },
      {
        id: '2',
        clientName: 'Mary Johnson',
        shiftTime: '14:00 - 17:00',
        location: '456 Oak Ave',
        status: 'completed',
        tasks: [],
        createdAt: new Date(),
        updatedAt: new Date(),
      },
    ]);

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText('John Smith')).toBeInTheDocument();
      expect(screen.getByText('Mary Johnson')).toBeInTheDocument();
    });

    // Filter by completed
    const select = screen.getByRole('combobox');
    fireEvent.change(select, { target: { value: 'completed' } });

    await waitFor(() => {
      expect(screen.queryByText('John Smith')).not.toBeInTheDocument();
      expect(screen.getByText('Mary Johnson')).toBeInTheDocument();
    });
  });

  it('should show "No schedules found" when no schedules exist', async () => {
    mockedMockApi.getSchedules.mockResolvedValue([]);
    mockedMockApi.getStats.mockResolvedValue({
      total: 0,
      missed: 0,
      upcoming: 0,
      completed: 0,
    });

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText('No schedules found')).toBeInTheDocument();
    });
  });

  it('should handle API errors gracefully', async () => {
    mockedMockApi.getSchedules.mockRejectedValue(new Error('API Error'));
    mockedMockApi.getStats.mockRejectedValue(new Error('API Error'));

    render(<Home />);

    await waitFor(() => {
      expect(screen.getByText('No schedules found')).toBeInTheDocument();
    });
  });
});