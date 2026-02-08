import { useState } from "react";
import api from "../api/api";
import { useAuth } from "../hooks/useAuth";
import { useNavigate, Link } from "react-router-dom";
import { jwtDecode } from "jwt-decode"; // Named import


const Login = () => {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            const response = await api.post("/users/login", { email, password });
            const { access_token } = response.data;
            const decoded = jwtDecode(access_token);
            const role = decoded.role || 'user';
            login(access_token, role);
            navigate("/rewardpage");
        } catch (error) {
            console.error("Login failed:", error);
            alert("Login failed. Please check your credentials.");
        }
    };

    return (
        <form onSubmit={handleSubmit} className="max-w-md mx-auto mt-20 p-6 bg-purple-50 rounded-lg shadow-md ">
            <h2 className="text-2xl font-bold mb-4">Login</h2>
            <div className="mb-4">
                <label htmlFor="email" className="block text-gray-700">Email:</label>
                <input
                    type="email"
                    id="email"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md"
                />
            </div>
            <div className="mb-4">
                <label htmlFor="password" className="block text-gray-700">Password:</label>
                <input
                    type="password"
                    id="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    className="w-full px-3 py-2 border rounded-md"
                />
            </div>
            <button type="submit" className="w-full bg-purple-500 text-white py-2 px-4 rounded-md hover:bg-purple-600">
                Login
            </button>
            <p className="mt-4 text-center">
                Don't have an account? <Link to="/register" className="text-purple-500 hover:underline">Register</Link>
            </p>
        </form>
    );
};

export default Login;

// AllowedOrigins: []string{"http://localhost:3000"}