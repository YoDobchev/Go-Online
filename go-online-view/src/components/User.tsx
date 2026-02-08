import React, { useCallback, useState } from "react";
import Navbar from "../components/Navbar";
import GameList from "../components/GameList";
import BottomCenMsg from "./BottomCenMsg";
import "../styles/Profile.scss";

type User = {
    email: string;
    username: string;
    elo: number;
    rank: string | number;
    role: string;
};

type ProfileProps = {
    user: User | null;
    loading?: boolean;
    canBan?: boolean
};

const User: React.FC<ProfileProps> = ({
    user,
    loading = false,
    canBan = false,
}) => {

    const [show, setShow] = useState(false);
    const [msg, setMsg] = useState("");
    const handleCloseMsg = useCallback(() => setShow(false), []);
    const handleBan = async () => {
        if (!user) return;
        try {
            await fetch(`/api/users/${user.username}`, {
                method: 'DELETE',
            });
        } catch {
            setMsg("Failed to ban user");
            setShow(true);
            return;
        } finally {
            setMsg("User has been banned");
            setShow(true);
        }
    };

    if (loading) return <div>Loading profile...</div>;

    if (!user) {
        return (
            <div className="must-login">
                You must be logged in to view this page.
            </div>
        );
    }

    const query = `q=${encodeURIComponent(user.username)}&status=finished`;

    return (
        <div>
            <Navbar />
            <div className="profile-cont">
                <div className="profile-stats">
                    <h1>Profile</h1>
                    <p>Email: {user.email}</p>
                    <p>Username: {user.username}</p>
                    <p>Elo: {user.elo}</p>
                    <p>Rank: {user.rank}</p>
                    <p>Role: {user.role}</p>

                    {canBan && <button className="banbtn" onClick={handleBan}>BAN</button>}
                </div>

                <div className="vertical"></div>

                <div className="profile-gamelist">
                    <h1>Past games</h1>
                    <GameList renderFilters={false} query={query} />
                </div>
            </div>

              <BottomCenMsg
                visible={show}
                message={msg}
                backgroundColor="#e90e0e"
                textColor="#f7f7f7"
                timeAfterFadeMs={1000}
                onClose={handleCloseMsg}
            />

        </div>
    );
};

export default User;
