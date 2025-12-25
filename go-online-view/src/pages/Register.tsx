import React, { useState } from "react";
import "../styles/auth.scss";
import { API_BASE } from "../config";

type RegisterResponse = {
    message?: string;
};

const Register: React.FC = () => {
    const [email, setEmail] = useState<string>("");
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [message, setMessage] = useState<string>("");

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();

        try {
            const res = await fetch(`${API_BASE}/auth/register`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({email, username, password }),
            });

            if (!res.ok) {
                throw new Error("Registration failed");
            }

            const data = (await res.json()) as RegisterResponse;

            setMessage(data.message ?? "no msg");
        } catch (err) {
            console.error(err);
            setMessage("failed or user already exists");
        }
    };

    return (
        <div className="auth">
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="email">Email</label>
                    <input
                        id="email"
                        type="email"
                        name="email"
                        value={email}
                        onChange={(e) => setEmail(e.target.value)}
                        required
                    />
                </div>

                <div>
                    <label htmlFor="username">Username</label>
                    <input
                        id="username"
                        type="text"
                        name="username"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
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

                <button type="submit">Register</button>
            </form>

            {message && <p>{message}</p>}
        </div>
    );
};

export default Register;
