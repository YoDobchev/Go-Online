import { useEffect, useState } from "react";
import "../styles/UsersDashboard.scss";
<<<<<<< Updated upstream
import Navbar from "../components/Navbar";

type User = {
    id: number;
    username: string;
    email: string;
    elo: number;
    rank: string;
    role: string;
};

export default function UsersDashboard() {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        fetch("/api/users", {
            credentials: "include",
        })
            .then((res) => {
                if (!res.ok) throw new Error("Failed to load users");
                return res.json();
            })
            .then(setUsers)
            .catch((err) => setError(err.message))
            .finally(() => setLoading(false));
    }, []);

    return (
        <>
            <Navbar />
            <div className="userlist">
                <div className="userlist-bar">
                    <div className="title">Users</div>
                </div>

                <div className="userlist-box">
                    <div className="userlist-header">
                        <div>ID</div>
                        <div>Username</div>
                        <div>Email</div>
                        <div>Elo</div>
                        <div>Rank</div>
                        <div>Role</div>
                    </div>

                    {loading && <div className="state">Loading users…</div>}
                    {error && <div className="state error">{error}</div>}
                    {!loading && users.length === 0 && (
                        <div className="state">No users found</div>
                    )}

                    {!loading &&
                        users.map((u) => (
                            <div key={u.id} className="userrow">
                                <div className="cell muted">{u.id}</div>
                                <div className="cell">{u.username}</div>
                                <div className="cell">{u.email}</div>
                                <div className="cell">{u.elo}</div>
                                <div className="cell">{u.rank}</div>
                                <div className="cell muted">{u.role}</div>
                            </div>
                        ))}
                </div>
            </div>
        </>
    );
=======

type User = {
  id: number;
  username: string;
  email: string;
  elo: number;
  rank: string;
  role: string;
};

export default function UsersDashboard() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch("/api/users", {
      credentials: "include",
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to load users");
        return res.json();
      })
      .then(setUsers)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="userlist">
      <div className="userlist-bar">
        <div className="title">Users</div>
      </div>

      <div className="userlist-box">
        <div className="userlist-header">
          <div>ID</div>
          <div>Username</div>
          <div>Email</div>
          <div>Elo</div>
          <div>Rank</div>
          <div>Role</div>
        </div>

        {loading && <div className="state">Loading users…</div>}
        {error && <div className="state error">{error}</div>}
        {!loading && users.length === 0 && (
          <div className="state">No users found</div>
        )}

        {!loading &&
          users.map((u) => (
            <div key={u.id} className="userrow">
              <div className="cell muted">{u.id}</div>
              <div className="cell">{u.username}</div>
              <div className="cell">{u.email}</div>
              <div className="cell">{u.elo}</div>
              <div className="cell">{u.rank}</div>
              <div className="cell muted">{u.role}</div>
            </div>
          ))}
      </div>
    </div>
  );
>>>>>>> Stashed changes
}
