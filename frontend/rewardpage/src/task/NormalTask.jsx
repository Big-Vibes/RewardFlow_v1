// NormalTasks component - Daily Task Checklist System
// CHANGE: Fixed debugging issues and added safety checks

import { completeTaskDaily, getDailyTasks } from "../api/apitask";
import { useState, useEffect } from "react";
import { CheckCircle2, Clock, AlertCircle } from "lucide-react";

const TASK_COUNT = 5;
const COOLDOWN_MINUTES = 5;
const POINTS_PER_TASK = 20;

function NormalTasks({ onComplete, onPointsUpdate }) {
  
  // CHANGE: State for 5-task daily checklist system
  const [tasks, setTasks] = useState([]);
  const [completedCount, setCompletedCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [loadingTaskId, setLoadingTaskId] = useState(null);
  const [error, setError] = useState(null);
  const [cooldownTime, setCooldownTime] = useState(null);
  const [isCooldownActive, setIsCooldownActive] = useState(false);

  // CHANGE: Load daily tasks from backend on component mount
  useEffect(() => {
    const loadDailyTasks = async () => {
      try {
        setLoading(true);
        // CHANGE: Fetch today's task state from backend
        // Backend endpoint: GET /api/tasks/daily
        // Returns: { tasks: [{ id, number, completed, completedAt }], completedCount, lastCompletedAt }
        const response = await getDailyTasks();
        
        // CHANGE: Added safety check for tasks array
        if (response && response.tasks && Array.isArray(response.tasks)) {
          setTasks(response.tasks);
          setCompletedCount(response.completedCount || 0);
          
          // CHANGE: If last task was just completed, set cooldown timer
          if (response.lastCompletedAt) {
            const lastTime = new Date(response.lastCompletedAt).getTime();
            const nextTime = lastTime + (COOLDOWN_MINUTES * 60 * 1000);
            
            // CHANGE: Start cooldown countdown if within cooldown window
            const now = Date.now();
            if (now < nextTime) {
              setIsCooldownActive(true);
              const remainingMs = nextTime - now;
              setCooldownTime(Math.ceil(remainingMs / 1000));
            }
          }
        } else {
          // CHANGE: Handle case where tasks array is missing or invalid
          console.warn('Invalid response format from getDailyTasks:', response);
          setError('Invalid task data received. Please refresh.');
          setTasks([]);
          setCompletedCount(0);
        }
      } catch (err) {
        console.error('Failed to load daily tasks:', err);
        // CHANGE: Better error message handling
        setError(err.message || 'Failed to load tasks. Please refresh.');
        setTasks([]);
      } finally {
        setLoading(false);
      }
    };

    loadDailyTasks();
  }, []);

  // CHANGE: Cooldown countdown timer - decrements every second
  useEffect(() => {
    // CHANGE: Added null check for cooldownTime
    if (!isCooldownActive || cooldownTime === null || cooldownTime <= 0) {
      setIsCooldownActive(false);
      setCooldownTime(null);
      return;
    }

    const interval = setInterval(() => {
      setCooldownTime((prev) => {
        // CHANGE: Added null check before decrementing
        if (prev === null) return null;
        const newTime = prev - 1;
        if (newTime <= 0) {
          setIsCooldownActive(false);
          return null;
        }
        return newTime;
      });
    }, 1000);

    return () => clearInterval(interval);
  }, [isCooldownActive, cooldownTime]);

  // CHANGE: Handle task completion with cooldown validation
  const handleTaskClick = async (taskId) => {
    // CHANGE: Added validation for taskId
    if (!taskId) {
      setError('Invalid task ID. Please try again.');
      return;
    }

    // CHANGE: Ensure taskId is a string (convert if it's an object)
    const taskIdStr = typeof taskId === 'string' ? taskId : taskId?.toString?.() || String(taskId);
    console.log('Attempting to complete task:', { taskId: taskIdStr, type: typeof taskId });

    // CHANGE: Prevent completion if daily limit reached
    if (completedCount >= TASK_COUNT) {
      setError('✅ All 5 daily tasks completed! Return tomorrow for more.');
      return;
    }

    // CHANGE: Prevent completion if cooldown is active (backend enforces this)
    if (isCooldownActive) {
      setError(`⏳ Please wait ${cooldownTime}s before completing another task (5-minute cooldown)`);
      return;
    }

    try {
      setLoadingTaskId(taskIdStr);
      setError(null);

      // CHANGE: Call backend API to complete task
      // Backend endpoint: POST /api/tasks/complete
      // Body: { taskId }
      // Backend validation:
      // - Checks if user already completed 5 tasks today
      // - Checks if user is within cooldown period (5 minutes)
      // - Updates lastCompletedAt timestamp
      // - Returns updated task state
      const response = await completeTaskDaily(taskIdStr);
      console.log('Task completion response:', response);

      // CHANGE: Added safety check for response
      if (response && !response.error && response.success) {
        // CHANGE: Update local state after successful backend completion
        setTasks((prev) =>
          prev.map((task) =>
            (task.id === taskIdStr)
              ? { ...task, completed: true, completedAt: new Date() }
              : task
          )
        );

        const newCompletedCount = completedCount + 1;
        setCompletedCount(newCompletedCount);

        // CHANGE: Set cooldown timer for next task (5 minutes)
        setIsCooldownActive(true);
        setCooldownTime(COOLDOWN_MINUTES * 60);

        // CHANGE: Calculate and update total points
        const totalPoints = newCompletedCount * POINTS_PER_TASK;
        if (typeof onPointsUpdate === 'function') {
          onPointsUpdate(totalPoints);
        }

        if (typeof onComplete === 'function') {
          onComplete(taskIdStr);
        }

        // CHANGE: Show completion message
        if (newCompletedCount === TASK_COUNT) {
          setError('🎉 All 5 tasks completed today! Great job!');
        }
      } else if (response?.error) {
        // CHANGE: Handle specific error responses from backend
        const errorMsg = response.error || 'Failed to complete task';
        setError(errorMsg.includes('cooldown') ? `⏳ ${errorMsg}` : errorMsg);
      } else {
        // CHANGE: Handle unsuccessful response
        setError(response?.message || 'Failed to complete task');
      }
    } catch (err) {
      console.error('Failed to complete task:', err);
      // CHANGE: Better error handling with fallback messages
      let errorMessage = 'Failed to complete task. Please try again.';
      
      if (err.response?.status === 429) {
        // CHANGE: Handle rate limiting (cooldown)
        errorMessage = `⏳ Cooldown active. Please wait before trying again.`;
      } else if (err.response?.data?.error) {
        errorMessage = err.response.data.error;
      } else if (err.response?.data?.message) {
        errorMessage = err.response.data.message;
      } else if (err.message) {
        errorMessage = err.message;
      }
      
      setError(errorMessage);
    } finally {
      setLoadingTaskId(null);
    }
  };

  // CHANGE: Format cooldown time as MM:SS with null check and better display
  const formatTime = (seconds) => {
    if (seconds === null || seconds === undefined || seconds <= 0) return '0:00';
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  // CHANGE: Calculate remaining tasks and points
  const totalPoints = completedCount * POINTS_PER_TASK;
  const progressPercent = tasks.length > 0 ? (completedCount / TASK_COUNT) * 100 : 0;
  const remainingTasks = TASK_COUNT - completedCount;

  if (loading) {
    return (
      <div className="w-full py-8 text-center">
        <p className="text-gray-600">⏳ Loading daily tasks...</p>
      </div>
    );
  }

  // CHANGE: Handle case where no tasks are loaded
  if (!tasks || tasks.length === 0) {
    return (
      <div className="w-full py-8 text-center">
        <p className="text-red-600">❌ No tasks available. Please refresh the page.</p>
      </div>
    );
  }

  return (
    <div className="w-full space-y-4">
      {/* Header with Points and Progress */}
      <div className="rounded-xl border border-gray-200 bg-white p-6 shadow-sm">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-lg font-semibold text-gray-900">Daily Task Checklist</h1>
          <div className="flex items-center gap-2 bg-purple-100 px-3 py-1 rounded-full">
            <span className="text-xl font-bold text-purple-600">{totalPoints}</span>
            <span className="text-sm text-purple-600">points today</span>
          </div>
        </div>

        {/* Progress Bar */}
        <div className="space-y-2">
          <div className="flex justify-between items-center">
            <span className="text-sm font-medium text-gray-700">
              {completedCount} of {TASK_COUNT} tasks completed
            </span>
            <span className="text-xs font-semibold text-gray-500">
              {remainingTasks} remaining • {Math.round(progressPercent)}%
            </span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3 overflow-hidden">
            <div
              className="bg-gradient-to-r from-purple-500 to-indigo-600 h-full transition-all duration-300"
              style={{ width: `${progressPercent}%` }}
            />
          </div>
        </div>

        {/* Cooldown Status */}
        {isCooldownActive && cooldownTime !== null && (
          <div className="mt-3 p-3 bg-yellow-50 border-2 border-yellow-200 rounded-lg flex items-center gap-3">
            <Clock className="h-5 w-5 text-yellow-600 flex-shrink-0" />
            <div className="flex-1">
              <span className="text-sm text-yellow-900 font-semibold">
                ⏳ Cooldown Active
              </span>
              <p className="text-xs text-yellow-700 mt-1">
                Next task available in <span className="font-mono font-bold text-yellow-600">{formatTime(cooldownTime)}</span>
              </p>
            </div>
          </div>
        )}

        {/* Error/Success Messages */}
        {error && (
          <div className={`mt-3 p-3 rounded-lg text-sm flex items-center gap-2 ${
            error.includes('🎉') || error.includes('✅')
              ? 'bg-green-50 border border-green-200 text-green-700'
              : 'bg-red-50 border border-red-200 text-red-700'
          }`}>
            {error.includes('🎉') || error.includes('✅') ? 
              <CheckCircle2 className="h-4 w-4" /> : 
              <AlertCircle className="h-4 w-4" />
            }
            {error}
          </div>
        )}
      </div>

      {/* 5 Task Buttons Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-2">
        {tasks.map((task, index) => {
          if (!task || typeof task !== "object") {
            return (
              <div
                key={index}
                className="h-20 rounded-md border border-red-200 text-xs flex items-center justify-center"
              >
                Invalid task
              </div>
            )
          }

          return (
            <div
              key={task.id || index}
              className={`rounded-md border transition-all duration-200 flex items-center justify-center h-20 ${
                task.completed
                  ? "border-purple-300 bg-purple-50"
                  : "border-gray-200 bg-white hover:border-purple-300 hover:shadow-sm"
              }`}
            >
              {task.completed ? (
                <div className="flex flex-col items-center gap-1">
                  <CheckCircle2 className="h-5 w-5 text-purple-600" />
                  <span className="text-[10px] font-semibold text-purple-700">
                    +{POINTS_PER_TASK} pts
                  </span>
                </div>
              ) : (
                <button
                  onClick={() => handleTaskClick(task.id)}
                  disabled={
                    loadingTaskId === task.id ||
                    isCooldownActive ||
                    completedCount >= TASK_COUNT
                  }
                  className={`w-full h-full flex flex-col items-center justify-center gap-1 rounded-md text-xs font-semibold transition-all ${
                    completedCount >= TASK_COUNT
                      ? "bg-gray-100 text-gray-400 cursor-not-allowed"
                      : isCooldownActive
                      ? "bg-purple-100 text-purple-600 cursor-wait"
                      : loadingTaskId === task.id
                      ? "bg-purple-100 text-purple-600 animate-pulse"
                      : "bg-purple-600 text-white hover:bg-purple-500 active:scale-95"
                  }`}
                  title={
                    isCooldownActive
                      ? `Wait ${cooldownTime}s before next task`
                      : "Complete task"
                  }
                >
                  <span>Task {task.number || index + 1}</span>
                  <span className="text-[10px]">+{POINTS_PER_TASK}</span>
                </button>
              )}
            </div>
          )
        })}
      </div>


      {/* Info Banner */}
      <div className="p-4 bg-purple-50 border border-gray-200 rounded-lg space-y-2">
        <p className="text-sm font-semibold text-purple-900">
          📋 Daily Task Checklist Rules
        </p>
        <ul className="text-xs text-purple-800 space-y-1 ml-4">
          <li>✓ Complete exactly <strong>5 tasks per day</strong></li>
          <li>✓ <strong>5-minute cooldown</strong> enforced between task completions</li>
          <li>✓ Earn <strong>{POINTS_PER_TASK} points per task</strong> ({TASK_COUNT * POINTS_PER_TASK} total per day)</li>
          <li>✓ Completed tasks become <strong className="text-green-700">disabled and highlighted</strong></li>
          <li>✓ Tasks <strong>unlock sequentially every 5 minutes</strong> after each completion</li>
          <li>✓ Tasks <strong>automatically reset at midnight</strong></li>
          {/* <li>✓ Completion history and progress <strong>saved in MongoDB</strong></li> */}
        </ul>
      </div>
    </div>
  );
}

export default NormalTasks;

{/* 5 Task Buttons Grid
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3">
        {tasks.map((task, index) => {
          // CHANGE: Added safety checks for task properties
          if (!task || typeof task !== 'object') {
            return <div key={index}>Invalid task</div>;
          }

          // CHANGE: Calculate margin increase for completed tasks (3x visual progression)
          // When completed, the task expands with increased margin for visual celebration
          const baseMargin = 0;
          const progressMargin = task.completed ? 24 : baseMargin; // 3x base unit increase
          const basePadding = 16; // p-4 = 1rem = 16px
          const progressPadding = task.completed ? basePadding * 3 : basePadding; // 3x padding expansion

          return (
            <div
              key={task.id || index}
              style={{ 
                margin: `${progressMargin}px`,
                padding: `${progressPadding}px`,
                transition: 'all 300ms cubic-bezier(0.4, 0, 0.2, 1)'
                
                
              }}
              className={`rounded-lg border-2 transition-all duration-300 flex flex-col items-center justify-center min-h-32 ${
                task.completed
                  ? 'border-green-300 bg-green-50 shadow-md'
                  : 'border-gray-200 bg-white hover:border-indigo-300 hover:shadow-md'
              }`}
            >
              {task.completed ? (
                // CHANGE: Show checkmark for completed tasks with celebration styling
                <div className="flex flex-col items-center gap-2">
                  <CheckCircle2 className="h-8 w-8 text-green-600 animate-bounce" />
                  <span className="text-xs font-bold text-green-700">COMPLETED</span>
                  <span className="text-xs text-green-600">+{POINTS_PER_TASK} pts</span>
                </div>
              ) : (
                <button
                  onClick={() => handleTaskClick(task.id)}
                  disabled={loadingTaskId === task.id || isCooldownActive || completedCount >= TASK_COUNT}
                  className={`w-full h-full flex flex-col items-center justify-center gap-2 rounded-md transition-all font-semibold ${
                    completedCount >= TASK_COUNT
                      ? 'bg-gray-100 text-gray-400 cursor-not-allowed'
                      : isCooldownActive
                      ? 'bg-yellow-100 text-yellow-600 cursor-wait'
                      : loadingTaskId === task.id
                      ? 'bg-indigo-100 text-indigo-600 animate-pulse'
                      : 'bg-indigo-600 text-white hover:bg-indigo-500 active:scale-95'
                  }`}
                  title={isCooldownActive ? `Wait ${cooldownTime}s before next task` : 'Complete this task'}
                >
                  <span className="text-sm font-bold">Task {task.number || index + 1}</span>
                  {loadingTaskId === task.id ? (
                    <span className="animate-spin text-lg">⏳</span>
                  ) : (
                    <span className="text-xs font-semibold">+{POINTS_PER_TASK}</span>
                  )}
                </button>
              )}
            </div>
          );
        })}
      </div> */}