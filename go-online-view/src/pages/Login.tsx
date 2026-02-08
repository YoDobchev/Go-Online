import React, { useState, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { API_BASE } from "../config";
import { UserContext } from "../context/UserContext";

const Login: React.FC = () => {
    const { setUser } = useContext(UserContext);
    const [identifier, setIdentifier] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [message, setMessage] = useState<string>("");
    const [loading, setLoading] = useState(false);

    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setLoading(true);
        try {
            const res = await fetch(`${API_BASE}/auth/login`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                credentials: "include",
                body: JSON.stringify({ identifier, password }),
            });

            if (!res.ok) throw new Error("Login failed");

            const meRes = await fetch(`${API_BASE}/auth/me`, { credentials: "include" });
            const me = meRes.ok ? await meRes.json() : null;

            setUser(me);
            setMessage("Logged in successfully");

            setTimeout(() => navigate("/"), 500);
        } catch (err) {
            console.error(err);
            setMessage("Invalid identifier or password");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="modal-backdrop">
            <div className="modal" role="dialog" aria-modal="true" aria-label="Login">
                <a
                    href="/"
                    style={{
                        display: "inline-block",
                        marginBottom: "0.75rem",
                        color: "#666",
                        textDecoration: "underline",
                        fontSize: "0.9rem",
                    }}
                >
                    Back
                </a>
                <div className="modal-title">Log in</div>

                <form onSubmit={handleSubmit}>
                    <div className="modal-row">
                        <label htmlFor="identifier">Email or Username</label>
                        <input
                            id="identifier"
                            type="text"
                            value={identifier}
                            onChange={(e) => setIdentifier(e.target.value)}
                            disabled={loading}
                            placeholder="Enter email or username"
                            required
                        />
                    </div>

                    <div className="modal-row">
                        <label htmlFor="password">Password</label>
                        <input
                            id="password"
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            disabled={loading}
                            placeholder="Enter your password"
                            required
                        />
                    </div>

                    {message && <div className="state error">{message}</div>}

                    <div className="modal-actions">
                        <button type="submit" disabled={loading}>
                            {loading ? "Logging in…" : "Log in"}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Login;