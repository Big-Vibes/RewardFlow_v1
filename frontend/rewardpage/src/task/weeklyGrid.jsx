// weeklyGrid (DailyStreak) component
// - `streak` expected shape: { mon: boolean, tue: boolean, ... } (lowercased keys)
// - `onCheckIn`: callback to register today's check-in
//
// Integration note:
// Frontend should POST to `/api/streak/update` when user checks in.
// The backend should return the updated streak object which the frontend stores.
import { checkIn } from "../api/apitask";
import { useState } from "react";
import { Calendar } from "lucide-react";

const days = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]

function DailyStreak({ streak, onCheckIn }) {
  const [loading, setLoading] = useState(false);

  const handleCheckIn = async () => {
    try {
      setLoading(true);
      await checkIn(onCheckIn); // checkIn will call setter passed in apitask
    } catch (err) {
      console.error('Check-in failed', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <>
    <div className="w-full space-y-4 rounded-xl border border-gray-200 bg-white shadow-sm ">
      <div className=" w-full text-lg font-semibold text-gray-800 px-6 py-4 
      bg-gray-50 rounded-tr-xl rounded-tl-xl border-b border-gray-200">

        <div className="flex items-center ">
          <Calendar className="inline h-5 w-5 mr-2 text-cyan-400" />
          <h2>Daily Streak</h2> 
        </div>
      </div>
    <div className="task-box p-6 ">
    

      <div className="grid grid-cols-7 gap-2">
        {days.map(day => (
          <div
            key={day}
            className={`px-2 py-2 text-center text-sm font-medium rounded-full border border-gray-300 shadow-sm ${
              streak[day.toLowerCase()] ? "bg-gray-500 text-white" : "bg-gray-200 text-gray-700"
            }`}
          >
            {day}
          </div>
          
        ))}
        <div className="col-span-7 mt-4 text-center text-sm text-gray-600">
          <span>
            click to win earn +5 points!
          </span>
        </div>
      </div>

      <button
        onClick={handleCheckIn}
        disabled={loading}
        className="mt-3 w-full bg-gray-600 text-white rounded-lg px-4 py-2 font-semibold hover:bg-gray-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
      >
        {loading ? 'Checking in...' : 'Check in Today'}
      </button>
    </div>
    </div>
  
    </>
  )
}

export default DailyStreak;
