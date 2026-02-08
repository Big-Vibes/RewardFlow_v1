import { useState } from "react"
import { StarIcon } from "@heroicons/react/24/solid"

export default function Example() {
  const POINTS_PER_TASK = 20
  const MAX_POINTS = 5000

  const [completedTasks, setCompletedTasks] = useState([false])

  const completedCount = completedTasks.filter(Boolean).length
  const totalPoints = completedCount * POINTS_PER_TASK

  const handleTaskClick = (index) => {
    setCompletedTasks((prev) =>
      prev.map((task, i) => (i === index ? true : task))
    )
  }

  return (
    <div className="flex gap-6">
      {completedTasks.map((isDone, index) => (
        <div
          key={index}
          className="w-1/3 rounded-xl border border-gray-200 bg-white p-5 shadow-sm"
        >
          {/* Header */}
          <div className="flex items-center justify-between mb-3">
            <h1 className="text-sm font-semibold text-gray-700">
              Points Balance
            </h1>

            <div className="flex items-center gap-1">
              <span className="text-lg font-bold text-indigo-600">
                {totalPoints}
              </span>
              <StarIcon className="h-5 w-5 text-yellow-500" />
            </div>
          </div>

          {/* Progress Info */}
          <div className="mb-3">
            <h1 className="text-sm font-medium text-gray-700">
              Progress to $5 gift card
            </h1>
            <p className="text-xs text-gray-500">
              {totalPoints}/{MAX_POINTS}
            </p>

            <progress
              className="w-full h-2 mt-2"
              value={totalPoints}
              max={MAX_POINTS}
            />
          </div>

          {/* Status */}
          <h1 className="text-xs text-gray-500 mb-3">
            {totalPoints === 0
              ? "Just getting started"
              : `${completedCount} task(s) completed`}
          </h1>

          {/* Action */}
          <button
            onClick={() => handleTaskClick(index)}
            disabled={isDone}
            className={`w-full rounded-md px-3 py-2 text-sm font-semibold text-white
              ${
                isDone
                  ? "bg-gray-400 cursor-not-allowed"
                  : "bg-indigo-600 hover:bg-indigo-500"
              }`}
          >
            {isDone ? "Task Completed" : "+20 Complete Task"}
          </button>
        </div>
      ))}
    </div>
  )
}
