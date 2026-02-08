import { useState, useEffect } from 'react';
import NormalTasks from '../task/NormalTask';
import DailyStreak from '../task/weeklyGrid';
import LeaderboardBox from '../task/leaderboard';
import api from '../api/api';
import { useAuth } from '../hooks/useAuth';
import { Trophy, TrendingUp } from 'lucide-react';

/**
 * TaskDashboard Component
 * 
 * Parent component that manages the three-task system:
 * 1. NormalTasks - daily task checklist
 * 2. DailyStreak - weekly check-in grid
 * 3. LeaderboardBox - user rankings
 * 
 * State Management:
 * - tasks: array of task objects { id, title, completed }
 * - streak: object { mon, tue, wed, thu, fri, sat, sun }
 * - leaderboard: array of ranked users
 * - loading: flag for API calls
 * - error: error message display
 * 
 * Backend Integration:
 * - Fetches tasks from GET /api/tasks
 * - Fetches streak from GET /api/streak
 * - Calls via completeTask, checkIn, loadLeaderboard helpers (apitask.ts)
 */
export default function TaskDashboard() {
  const { user } = useAuth();
  const [tasks, setTasks] = useState([]);
  const [streak, setStreak] = useState({});
  const [leaderboard, setLeaderboard] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  // Fetch leaderboard data with real-time polling
  const fetchLeaderboard = async () => {
    try {
      const response = await api.get('/api/leaderboard?limit=20');
      setLeaderboard(response.data || []);
    } catch (err) {
      console.error('Failed to fetch leaderboard:', err);
      // setError('Failed to load leaderboard. Please try again later.');
    }
  };

  // Fetch initial data
  useEffect(() => {
    const fetchInitialData = async () => {
      try {
        setError(null);
        const [tasksRes, streakRes] = await Promise.all([
          api.get('/tasks'),
          api.get('/streak'),
        ]);
        setTasks(tasksRes.data || []);
        setStreak(streakRes.data || {});
        
        // Fetch initial leaderboard
        await fetchLeaderboard();
      } catch (err) {
        console.error('Failed to fetch initial data:', err);
        setError('Failed to load tasks and streak. Please try again.');
      } finally {
        setLoading(false);
      }
    };

    fetchInitialData();
  }, []);

  // Real-time leaderboard polling (every 5 seconds)
  useEffect(() => {
    const interval = setInterval(async () => {
      await fetchLeaderboard();
    }, 5000);

    return () => clearInterval(interval);
  }, []);

  // Handle task completion - update local state
  const handleCompleteTask = (taskId) => {
    setTasks((prev) =>
      prev.map((task) =>
        task.id === taskId ? { ...task, completed: true } : task
      )
    );
  };

  // Handle streak update - state is updated by checkIn helper
  const handleCheckIn = (updatedStreak) => {
    setStreak(updatedStreak);
  };

  // Handle leaderboard open - refetch latest data
  const handleOpenLeaderboard = async () => {
    try {
      await fetchLeaderboard();
    } catch (err) {
      console.error('Failed to refresh leaderboard:', err);
    }
  };

  if (loading) {
    return (
      <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="text-center py-12">
          <p className="text-gray-500">Loading tasks...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
     <h1 className="mb-2 text-xl font-bold text-gray-900">
      Your reward journey
      </h1>
        
      {/* Error Banner */}
      {error && (
        <div className="mb-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}

      {/* Task Checklist, Streak, and Leaderboard Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Daily Tasks */}
        <div>
          <NormalTasks 
          tasks={tasks} 
          onComplete={handleCompleteTask} />
        </div>

        {/* Daily Streak */}
        <div>
          <DailyStreak 
          streak={streak} 
          onCheckIn={handleCheckIn} />
        </div>

        {/* Leaderboard */}
        <div>
          <LeaderboardBox onOpen={handleOpenLeaderboard} setLeaderboard={setLeaderboard} />
        </div>
      </div>

      {/* Leaderboard Results */}
      {leaderboard.length > 0 && (
        <div className="mt-8">
          <div className="rounded-xl border border-gray-200 bg-white shadow-sm overflow-hidden">
            {/* Header */}
            <div className="bg-gradient-to-r from-purple-600 to-indigo-600 px-6 py-4">
              <div className="flex items-center gap-2">
                <Trophy className="h-6 w-6 text-yellow-300" />
                <h2 className="text-2xl font-bold text-white">Leaderboard Rankings</h2>
              </div>
              <p className="text-purple-100 text-sm mt-1">Top 20 users by total points</p>
            </div>

            {/* Leaderboard Table */}
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead>
                  <tr className="border-b border-gray-200 bg-gray-50">
                    <th className="px-6 py-3 text-left">Rank</th>
                    <th className="px-6 py-3 text-left">Username</th>
                    <th className="px-6 py-3 text-right">Points</th>
                    <th className="px-6 py-3 text-center">Status</th>
                  </tr>
                </thead>
                <tbody>
                  {leaderboard.map((userRow, idx) => {
                    const isCurrentUser = user?.username === userRow.username;
                    const medalEmoji = idx === 0 ? '🥇' : idx === 1 ? '🥈' : idx === 2 ? '🥉' : '';
                    
                    return (
                      <tr
                        key={idx}
                        className={`border-b border-gray-100 transition-colors ${
                          isCurrentUser
                            ? 'bg-indigo-50 border-l-4 border-l-indigo-500 font-semibold'
                            : 'hover:bg-gray-50'
                        }`}
                      >
                        {/* Rank */}
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            <span className="text-lg font-bold text-gray-900">
                              {medalEmoji || `#${userRow.rank || idx + 1}`}
                            </span>
                          </div>
                        </td>

                        {/* Username */}
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            <span className="font-medium text-gray-900">{userRow.username || userRow.email}</span>
                            {isCurrentUser && (
                              <span className="inline-block px-2 py-1 bg-indigo-500 text-white text-xs font-bold rounded-full">
                                YOU
                              </span>
                            )}
                          </div>
                        </td>

                        {/* Points */}
                        <td className="px-6 py-4 text-right">
                          <div className="flex items-center justify-end gap-1">
                            <TrendingUp className={`h-4 w-4 ${isCurrentUser ? 'text-indigo-600' : 'text-gray-400'}`} />
                            <span className={`font-bold text-lg ${
                              isCurrentUser ? 'text-indigo-600' : 'text-gray-900'
                            }`}>
                              {userRow.points || 0}
                            </span>
                          </div>
                        </td>

                        {/* Status Badge */}
                        <td className="px-6 py-4 text-center">
                          {idx === 0 && <span className="inline-block px-3 py-1 bg-yellow-100 text-yellow-800 text-xs font-semibold rounded-full">Leader</span>}
                          {idx > 0 && idx <= 2 && <span className="inline-block px-3 py-1 bg-blue-100 text-blue-800 text-xs font-semibold rounded-full">Top 3</span>}
                          {idx > 2 && idx <= 9 && <span className="inline-block px-3 py-1 bg-green-100 text-green-800 text-xs font-semibold rounded-full">Top 10</span>}
                          {idx >= 10 && <span className="inline-block px-3 py-1 bg-gray-100 text-gray-800 text-xs font-semibold rounded-full">Rising</span>}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>

            {/* Footer */}
            <div className="px-6 py-4 bg-gray-50 border-t border-gray-200">
              <p className="text-xs text-gray-500">
                ✨ Updates every 5 seconds | {leaderboard.length} users ranked
              </p>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
