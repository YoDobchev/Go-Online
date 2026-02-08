import React, { useContext, useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import User from "../components/User";
import Navbar from "../components/Navbar";
import { UserContext } from "../context/UserContext";
import { API_BASE } from "../config";

type RemoteUser = {
    email: string;
    username: string;
    elo: number;
    rank: string | number;
    role: string;
};

const UserPage: React.FC = () => {
    const { user: me, loading: meLoading } = useContext(UserContext);
    const { username } = useParams<{ username: string }>();

    const [target, setTarget] = useState<RemoteUser | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const isModerator = !!me && (me.role === "admin" || me.role === "moderator");

    useEffect(() => {
        if (!isModerator) return;
        if (!username) return;

        const ctrl = new AbortController();
        (async () => {
            try {
                setLoading(true);
                setError(null);

                const res = await fetch(
                    `${API_BASE}/users/${encodeURIComponent(username)}`,
                    {
                        method: "GET",
                        credentials: "include",
                        signal: ctrl.signal,
                        headers: { Accept: "application/json" },
                    },
                );

                if (!res.ok) {
                    const txt = await res.text();
                    throw new Error(
                        txt || `Failed to load user (${res.status})`,
                    );
                }

                const data = (await res.json()) as Partial<RemoteUser>;
                setTarget({
                    email: data.email ?? "",
                    username: data.username ?? "",
                    elo: (data.elo as number) ?? 0,
                    rank: data.rank ?? "—",
                    role: data.role ?? "",
                });
            } catch (e: unknown) {
                if (e instanceof DOMException && e.name === "AbortError")
                    return;
                setError(
                    e instanceof Error ? e.message : "Something went wrong",
                );
            } finally {
                setLoading(false);
            }
        })();

        return () => ctrl.abort();
    }, [username, isModerator]);

    if (meLoading) {
        return (
            <>
                <Navbar />
                <div style={{ padding: 24 }}>Loading…</div>
            </>
        );
    }

    if (!isModerator) {
        return (
            <>
                <Navbar />
                <div className="flex flex-col items-center justify-center h-screen">
                    <h1 className="text-3xl font-bold mb-4">Access Denied</h1>
                    <p className="text-gray-600">
                        You do not have permission to view this page.
                    </p>
                </div>
            </>
        );
    }

    if (error) {
        return (
            <>
                <Navbar />
                <div style={{ padding: 24, color: "#b91c1c" }}>
                    Error: {error}
                </div>
            </>
        );
    }

    return <User user={target} loading={loading} />;
};

export default UserPage;
