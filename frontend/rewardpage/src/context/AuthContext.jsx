import { createContext, useEffect, useState } from "react";
import { jwtDecode } from "jwt-decode"; // Named import

const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [auth, setAuth] = useState({
        token: null,
        role: null,
        isAuthenticated: false,
        user: null,
    });

    useEffect(() => {
        const initializeAuth = async () => {
            const token = localStorage.getItem('token');
            if (token) {
                try {
                    const decoded = jwtDecode(token);
                    setAuth({
                        token,
                        role: decoded.role || 'user',
                        user: {
                            id: decoded.id || decoded.sub,
                            username: decoded.username,
                            email: decoded.email,
                        },
                        isAuthenticated: true,
                    });
                } catch (error) {
                    console.error('Invalid token:', error);
                    localStorage.removeItem('token');
                    localStorage.removeItem('role');
                }
            }
        };
        initializeAuth();
    }, []);

    const login = (token, role, userData = {}) => {
        localStorage.setItem('token', token);
        localStorage.setItem('role', role);
        setAuth({
            token,
            role,
            user: {
                id: userData.id,
                username: userData.username,
                email: userData.email,
            },
            isAuthenticated: true,
        });
    };

    const logout = () => {
        localStorage.removeItem('token');
        localStorage.removeItem('role');
        setAuth({
            token: null,
            role: null,
            user: null,
            isAuthenticated: false,
        });
    };

    return (
        <AuthContext.Provider value={{ ...auth, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

export { AuthContext };
