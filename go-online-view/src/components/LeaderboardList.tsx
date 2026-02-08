import { useEffect, useState } from "react";
import "../styles/Leaderboard.scss";

type LeaderboardUser = {
    id: number;
    username: string;
    elo: number;
    role: string;
    rank: string;
};

export default function LeaderboardList() {
    const [users, setUsers] = useState<LeaderboardUser[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetch("/api/leaderboard")
            .then((res) => {
                if (!res.ok) throw new Error("Failed to load leaderboard");
                return res.json();
            })
            .then(setUsers)
            .catch((err) => setError(err.message))
            .finally(() => setLoading(false));
    }, []);

    return (
        <div className="leaderboardlist">
            <div className="leaderboardlist-box">
                <div className="leaderboardlist-header">
                    <div>#</div>
                    <div>Username</div>
                    <div>Elo</div>
                    <div>Rank</div>
                </div>

                {loading && <div className="state">Loading leaderboard…</div>}
                {error && <div className="state error">{error}</div>}
                {!loading && !error && users.length === 0 && (
                    <div className="state">No users found</div>
                )}

                {!loading &&
                    !error &&
                    users.map((u, index) => (
                        <div key={u.id} className="leaderboardrow">
                            <div className="cell muted">{index + 1}</div>
                            <div className="cell">{u.username}</div>
                            <div className="cell">{u.elo}</div>
                            <div className="cell muted">{u.rank}</div>
                        </div>
                    ))}
            </div>
        </div>
    );
}