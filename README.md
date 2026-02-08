# RewardFlow_v1
RewardFlow is a full-stack reward and task-based engagement platform. The project demonstrates real-world system design, combining a modern React frontend with a Go backend and a scalable database architecture.
The RewardFlow platform allows users to complete tasks, build daily streaks, earn points, and compete on a leaderboard, all secured with authentication and role-based access control.
Features:
🔐 User Authentication (Register / Login / Logout)
🎯 Task Checklist System (+20 points per task)
🔁 Daily Task Reset Logic
📆 Weekly Streak Tracking (Monday – Sunday)
📊 Progress Bar & Reward Milestones
🏆 Leaderboard (Ranked by points)
🧩 Role-Based Authorization (RBAC)
⏳ Cooldown & Anti-abuse logic
🎨 Responsive UI with Tailwind CSS

Core Modules
1. Authentication: JWT-based login & registration, Password hashing with bcrypt, Protected routes using middleware.
2. Task Engine: Daily task checklist, Fixed reward per task, Cooldown enforcement, Auto-reset logic
3. Daily Streak System: Monday → Sunday tracking, One check-in per day, Streak persistence, Missed-day reset logic
4. Leaderboard: Global ranking, sorted by total points, Read-only API, Extendable to weekly/monthly boards

Technology Stack
- Frontend: React.js, Tailwind CSS, JavaScript / TypeScript, HeroIcons.
- Backend: Go (Golang), Gorilla Mux, JWT Authentication, bcrypt (password hashing).
- Database: MongoDB (Document-based).
- Dev Tools: Git & GitHub, VS Code, Thunder Client Extension(VS Code).

What You Learn From This Project:
* Backend: Rest API design in Go, Middleware patterns, JWT validation, Role-based access control, clean service & controller separation,  Database modeling (NoSQL), Time-based logic (daily reset, streaks).

* Frontend: Component-driven UI, State management patterns, Conditional rendering, Async API integration, UX feedback (loading, disabled states), Responsive layouts.

* System Design: Monorepo structuring, Separation of concerns, Security thinking, Scalable reward systems, Real-world business rules.
