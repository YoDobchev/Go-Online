import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { API_BASE } from "../config";

const Register: React.FC = () => {
    const [email, setEmail] = useState<string>("");
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [message, setMessage] = useState<string>("");
    const [loading, setLoading] = useState<boolean>(false);

    const navigate = useNavigate();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setLoading(true);

        try {
            const res = await fetch(`${API_BASE}/auth/register`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ email, username, password }),
            });

            if (!res.ok) throw new Error("Registration failed");

            const data = await res.json();
            setMessage(data.message ?? "Registration successful");

            setTimeout(() => navigate("/login"), 500);
        } catch (err) {
            console.error(err);
            setMessage("Failed or user already exists");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="modal-backdrop">
            <div
                className="modal"
                role="dialog"
                aria-modal="true"
                aria-label="Register"
            >
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

                <div className="modal-title">Register</div>

                <form onSubmit={handleSubmit}>
                    <div className="modal-row">
                        <label htmlFor="email">Email</label>
                        <input
                            id="email"
                            type="email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            disabled={loading}
                            placeholder="Enter your email"
                            required
                        />
                    </div>

                    <div className="modal-row">
                        <label htmlFor="username">Username</label>
                        <input
                            id="username"
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                            disabled={loading}
                            placeholder="Choose a username"
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
                            placeholder="Enter a password"
                            required
                        />
                    </div>

                    {message && <div className="state error">{message}</div>}

                    <div className="modal-actions">
                        <button type="submit" disabled={loading}>
                            {loading ? "Registering…" : "Register"}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};

export default Register;
