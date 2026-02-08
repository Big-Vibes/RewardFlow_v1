import { StarIcon, Share2Icon } from "lucide-react"

export default function EarnPoints() {
  return (
    <>
     <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
      <h1 className="mb-4 text-xl font-bold text-gray-900">
        Earn More Points
      </h1>
      
      <div className="flex flex-row gap-2">
        
        {/* Card 1 */}
        <div className="flex flex-col gap-2 rounded-lg bg-purple-50 p-4 border border-purple-100">
          <div className="flex items-center gap-3">
            <div className="flex items-center justify-center rounded-full bg-purple-200 p-2">
              <StarIcon className="h-4 w-4 text-purple-700" />
            </div>
            <h2 className="text-lg font-semibold text-gray-900">
              Refer and win 10,000 points
            </h2>
          </div>

          <p className="text-sm text-gray-600">
            Invite 5 friends by Feb 20 and earn a chance to be one of 5 winners of{" "}
            <span className="font-semibold text-purple-700">10,000 points</span>.
            Friends must complete onboarding with quality activity.
          </p>
        </div>

        {/* Card 2 */}
        <div className="flex flex-col gap-2 rounded-lg bg-purple-50 p-4 border border-purple-100">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="flex items-center justify-center rounded-full bg-purple-200 p-2">
                <Share2Icon className="h-4 w-4 text-purple-700" />
              </div>
              <h2 className="text-lg font-semibold text-gray-900">
                Share your stack
              </h2>
            </div>

            <span className="text-sm font-semibold text-purple-700">
              +25 pts
            </span>
          </div>

          <p className="text-sm text-gray-600">
            Share your true stack and earn bonus points.
          </p>
        </div>

      </div>
    </div>


    </>
  )
}



