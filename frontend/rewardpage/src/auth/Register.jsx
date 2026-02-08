import { useState } from "react";
import api from "../api/api";
import { useNavigate, Link } from "react-router-dom";

const Register = () => {
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const navigate = useNavigate();

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            await api.post("/users/register", { username, email, password });
            // Registration successful, redirect to login page or directly to reward page
            navigate("/rewardpage");
        } catch (error) {
            console.error("Registration failed:", error);
        }
    };

    return (
        <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
            <div className="sm:mx-auto sm:w-full sm:max-w-sm">
                <h2 className="mt-10 text-center text-2xl/9 font-bold tracking-tight text-purple-600">Sign up to your account</h2>
            </div>

            <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
                <form onSubmit={handleSubmit} className="space-y-6">
                    <div>
                        <label htmlFor="username" className="block text-sm/6 font-medium text-purple-600">
                            Username
                        </label>
                        <div className="mt-2">
                            <input
                                id="username"
                                name="username"
                                type="text"
                                required
                                value={username}
                                onChange={(e) => setUsername(e.target.value)}
                                autoComplete="username"
                                className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base text-black outline-2 -outline-offset-2 
                                outline-indigo-300 placeholder:text-purple-500 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6"
                            />
                        </div>
                        <label htmlFor="email" className="block text-sm/6 font-medium text-purple-600">
                            Email address
                        </label>
                        <div className="mt-2">
                            <input
                                id="email"
                                name="email"
                                type="email"
                                required
                                value={email}
                                onChange={(e) => setEmail(e.target.value)}
                                autoComplete="email"
                                className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base text-black outline-2 -outline-offset-2 outline-indigo-300 placeholder:text-blue-500 focus:outline-2 focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6"
                            />
                        </div>
                    </div>

                    <div>
                        <div className="flex items-center justify-between">
                            <label htmlFor="password" className="block text-sm/6 font-medium text-purple-600">
                                Password
                            </label>
                        </div>
                        <div className="mt-2">
                            <input
                                id="password"
                                name="password"
                                type="password"
                                required
                                value={password}
                                onChange={(e) => setPassword(e.target.value)}
                                autoComplete="current-password"
                                className="block w-full rounded-md bg-white/5 px-3 py-1.5 text-base text-black 
                                outline-2 -outline-offset-2 outline-indigo-300 placeholder:text-blue-500 focus:outline-2 
                                focus:-outline-offset-2 focus:outline-indigo-500 sm:text-sm/6"
                            />
                        </div>
                    </div>

                    <div>
                        <button
                            type="submit"
                            className="flex w-full justify-center rounded-md bg-purple-500 px-3 py-1.5 text-sm/6 font-semibold 
                            text-white hover:bg-purple-400 focus-visible:outline-2 focus-visible:outline-offset-2 
                            focus-visible:outline-purple-500">
                            Sign up
                        </button>
                    </div>
                </form>
                <p className="mt-10 text-center text-sm/6 text-gray-400">
                    Already have an account?{' '}
                    <Link to="/login" className="font-semibold text-indigo-400 hover:text-indigo-300">
                        Sign in
                    </Link>
                </p>
            </div>
        </div>
    );
};

export default Register;