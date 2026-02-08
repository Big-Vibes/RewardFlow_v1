import { PersonStandingIcon } from "lucide-react"

export default function ReferPoint(){
    return(
        <>
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="mb-4 text-xl font-bold text-gray-900">
            Refer and Earn
        </h1>

        <div className="flex flex-col gap-4 rounded-lg bg-purple-50 p-5 border border-purple-100">
            
            {/* Header */}
            <div className="flex items-start gap-3">
            <div className="flex h-9 w-9 items-center justify-center rounded-full bg-purple-200">
                <PersonStandingIcon className="h-4 w-4 text-purple-700" />
            </div>

            <div>
                <h2 className="text-base font-semibold text-gray-900">
                Share your link
                </h2>
                <p className="text-sm text-gray-600">
                Invite friends and earn <span className="font-medium text-purple-700">25 pts</span> when they join
                </p>
            </div>
            </div>

            {/* Stats */}
            <div className="flex justify-between rounded-md bg-white px-6 py-4 text-center">
            <div>
                <p className="text-lg font-semibold text-gray-900">0</p>
                <p className="text-xs text-gray-500">Referrals</p>
            </div>

            <div>
                <p className="text-lg font-semibold text-gray-700">0</p>
                <p className="text-xs text-gray-500">Points Earned</p>
            </div>
            </div>

        </div>
        </div>

        </>
    )
}