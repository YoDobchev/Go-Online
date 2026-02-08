import React, { useContext, useEffect, useRef, useState } from "react";
import { UserContext } from "../context/UserContext";
import "../styles/Navbar.scss";
import { API_BASE } from "../config";

const Navbar: React.FC = () => {
    const { user, loading, setUser } = useContext(UserContext);
    const [open, setOpen] = useState(false);
    const menuRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        const onClick = (e: MouseEvent) => {
            if (!menuRef.current) return;
            if (!menuRef.current.contains(e.target as Node)) setOpen(false);
        };
        const onKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") setOpen(false);
        };

        document.addEventListener("mousedown", onClick);
        document.addEventListener("keydown", onKeyDown);
        return () => {
            document.removeEventListener("mousedown", onClick);
            document.removeEventListener("keydown", onKeyDown);
        };
    }, []);

    const handleLogout = () => {
        fetch(`${API_BASE}/auth/logout`, {
            method: "DELETE",
            credentials: "include",
        });
        setUser(null);
        setOpen(false);
    };

    return (
        <nav className="navbar">
            <div className="nav-left">
                <a href="/" className="brand">
                    <img src="images/logo.jpg" alt="Go Online" />
                    <span>Go Online</span>
                </a>

                <a href="/search">Search</a>
                <a href="/blogs">Blogs</a>
                <a href="/leaderboard">LeaderBoard</a>
                <a href="/about">About</a>

            </div>

            <div className="nav-right">
                {loading ? (
                    <span className="loading">Loading…</span>
                ) : user ? (
                    <div className="profile-wrap" ref={menuRef}>
                        <button
                            className="icon-btn"
                            onClick={() => setOpen((v) => !v)}
                            aria-haspopup="menu"
                            aria-expanded={open}
                            aria-label="Open profile menu"
                        >
                            <svg viewBox="0 0 24 24" aria-hidden="true">
                                <path d="M12 12a4.5 4.5 0 1 0-4.5-4.5A4.5 4.5 0 0 0 12 12Zm0 2c-4.1 0-7.5 2.2-7.5 5v1h15v-1c0-2.8-3.4-5-7.5-5Z" />
                            </svg>
                        </button>

                        {open && (
                            <div className="menu" role="menu">
                                <a
                                    className="menu-item"
                                    href="/profile"
                                    role="menuitem"
                                    onClick={() => setOpen(false)}
                                >
                                    Profile
                                </a>
                                <button
                                    className="menu-item danger"
                                    role="menuitem"
                                    onClick={handleLogout}
                                >
                                    Logout
                                </button>
                            </div>
                        )}
                    </div>
                ) : (
                    <div className="auth">
                        <a href="/login">Login</a>
                        <a href="/register" className="register">
                            Register
                        </a>
                    </div>
                )}
            </div>
        </nav>
    );
};

export default Navbar;
