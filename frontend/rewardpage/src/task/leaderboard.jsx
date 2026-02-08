// Leaderboard Button Component - Triggers leaderboard refresh/fetch
// Integration: Connected to TaskDashboard real-time polling
import { useState } from "react";
import { Trophy } from "lucide-react";

function LeaderboardBox({ onOpen }) {
  const [loading, setLoading] = useState(false);

  const handleOpen = async () => {
    try {
      setLoading(true);
      if (typeof onOpen === 'function') {
        await onOpen();
      }
    } catch (err) {
      console.error('Failed to load leaderboard', err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="rounded-xl border border-gray-200 bg-gradient-to-br from-violet-500 to-fuchsia-500 p-6 shadow-sm">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <Trophy className="h-6 w-6 text-yellow-300" />
          <h2 className="text-lg font-bold text-white">Leaderboard</h2>
        </div>

        <button
          onClick={handleOpen}
          disabled={loading}
          className="px-4 py-2 bg-sky-900 text-white text-sm font-semibold rounded-lg hover:bg-sky-800 disabled:bg-gray-500 disabled:cursor-not-allowed transition-colors"
        >
          {loading ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
      
      <p className="text-sm text-white/80 mt-2">
        See who's at the top and climb the rankings!
      </p>
      <div className="p-4 bg-blue-50 border border-blue-200 rounded-lg mt-4">
        <p className="text-sm text-blue-700">
          💡 <strong>Leaderboard Tip:</strong> Complete up to 5 tasks daily to earn points. Climb the rankings and earn rewards!
        </p>
      </div>
    </div>

  );
}
export default LeaderboardBox;