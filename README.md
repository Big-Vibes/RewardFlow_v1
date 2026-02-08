# RewardFlow_v1
RewardFlow is a full-stack reward and task-based engagement platform. The project demonstrates real-world system design, combining a modern React  Go backend and a scalable database architecture.
The RewardFlow platform enables users to complete tasks, build daily streaks, earn points, and compete on a leaderboard, all while ensuring security through authentication and role-based access control.

**Features:**
- 🔐 User Authentication (Register / Login / Logout)
- 🎯 Task Checklist System (+20 points per task)
- 🔁 Daily Task Reset Logic
- 📆 Weekly Streak Tracking (Monday – Sunday)
- 📊 Progress Bar & Reward Milestones
- 🏆 Leaderboard (Ranked by points)
- 🧩 Role-Based Authorization (RBAC)
- ⏳ Cooldown & Anti-abuse logic
- 🎨 Responsive UI with Tailwind CSS

**Core Modules**
1. Authentication: JWT-based login & registration, Password hashing with bcrypt, Protected routes using middleware.
2. Task Engine: Daily task checklist, Fixed reward per task, Cooldown enforcement, Auto-reset logic
3. Daily Streak System: Monday → Sunday tracking, One check-in per day, Streak persistence, Missed-day reset logic
4. Leaderboard: Global ranking, sorted by total points, Read-only API, Extendable to weekly/monthly boards

**Technology Stack**
- Frontend: React.js, Tailwind CSS, JavaScript / TypeScript, HeroIcons.
- Backend: Go (Golang), Gorilla Mux, JWT Authentication, bcrypt (password hashing).
- Database: MongoDB (Document-based).
- Dev Tools: Git & GitHub, VS Code, Thunder Client Extension(VS Code).

**Architecture Overview**
* [ React Frontend ]
*         |
*         | REST API (JSON)
*         v
* [ Go API Server ]
*         |
*         v
* [ MongoDB Database ]

* Frontend handles UI, state, and API calls, 
* Backend manages authentication, business logic, streak rules, and leaderboard ranking. 
* MongoDB stores users, tasks, streaks, and points

**Project Structure:**
- Backend
*   backend/
*   ├── controllers/ 
*   ├── services/ 
*   ├── repositories/ 
*   ├── models/
*   ├── routes/ 
*   ├── utils/
*   ├── middleware/ 
*   └── main.go
* Frontend
* frontend/rewardpage
* ├── src/
* │   ├── api/ 
* │   ├── components/ 
* │   ├── body/ 
* │   ├── hooks/
* │   ├── context/ 
* │   ├── auth/
* │   ├── task/
* │   └── main.jsx

**Installation**
Prerequisites
* Node.js (>= 18)
* Go (>= 1.21)
* MongoDB (Cloud "Atlas")

**Clone the repository**
* [git clone https://github.com/your-username/rewardhub.git](https://github.com/Big-Vibes/RewardFlow_v1.git)
* cd rewardhub

**Running the Application**
* cd backend
* go mod tidy
* go run main.go
* go get go.mongodb.org/mongo-driver/v2/mongo
* go get -u github.com/gorilla/mux
* go get -u github.com/golang-jwt/jwt/v5

Backend runs on:
* http://Localhost:4000

* cd frontend/rewardpage
* npm install
* npm run dev
* npm create vite@latest my-project
* npm i lucide-react
* npm i @heroicons/react

frontend runs on:
* http/localhost:5173

**Configuration**
* Create a .env file in the backend directory:
* MONGO_URI=your_mongodb_Key
* JWT_SECRET=your_secret_key

**Development Roadmap**
* Phase 1(core, week1)
  Data collection,
  Backend connection,
  Database modeling (NoSQL),

* Phase 2 (Engagement, week 2-4)
  frontend connection,
  Auth (Login / Register / Logout),
  Normal task list,
  Daily streak logic,
  points system,
  leaderboard,

* Phase 3 (Polish)
  Progress indicators,
  Mobile responsiveness,
  Deployment on GitHub page, Vercel,




**What You Learn From This Project:**

Full‑stack architecture design, 
Authentication flows, 
Date‑based streak logic, 
REST API design, 
MongoDB schema modeling, 
Scalable frontend structure, 
Real‑world product thinking, 

📌 **Development Tip**

* This project is ideal as a public repository for recruiters it demonstrates logic, structure, and decision‑making beyond simple CRUD apps.
