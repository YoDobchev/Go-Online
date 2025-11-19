import React, { useState, useContext } from "react";
import "../styles/auth.scss";
import { useNavigate } from "react-router-dom";
import { API_BASE } from "../config";
import { UserContext } from "../context/UserContext";

const Login: React.FC = () => {
    const { setUser } = useContext(UserContext);
    const [identifier, setIdentifier] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [message, setMessage] = useState<string>("");
    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        try {
            const res = await fetch(`${API_BASE}/api/auth/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({ identifier, password }),
            });

            if (!res.ok) {
                throw new Error("Login failed");
            }

            const meRes = await fetch(`${API_BASE}/api/auth/me`, {
                credentials: "include",
            });
            const me = meRes.ok ? await meRes.json() : null;

            setUser(me);

            setMessage(`logged in`);

            setTimeout(() => navigate("/"), 500);
        } catch (err) {
            console.error(err);
            setMessage("invalid identifier or password");
        }
    };

    return (
        <div className="auth">
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="identifier">Email or Password</label>
                    <input
                        id="identifier"
                        type="text"
                        name="identifier"
                        value={identifier}
                        onChange={(e) => setIdentifier(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label htmlFor="password">Password</label>
                    <input
                        id="password"
                        type="password"
                        name="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        required
                    />
                </div>

                <button type="submit">Log in</button>
            </form>

            {message && <p>{message}</p>}
        </div>
    );
};

export default Login;
