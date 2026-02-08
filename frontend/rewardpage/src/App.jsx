
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Register from './auth/Register';
import Login from './auth/login';
import RewardPage from './body/rewardpage';
import ProtectedRoute from './component/ProtectedRoute';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/register" element={<Register />} />
        <Route path="/login" element={<Login />} />
        <Route
          path="/rewardpage"
          element={
            <ProtectedRoute>
              <RewardPage />
            </ProtectedRoute>
          }
        />
        <Route path="/" element={<Login />} />
      </Routes>
    </Router>
  );
}

export default App;
