import { useContext } from "react";
import Navbar from "../components/Navbar";
import GameList from "../components/GameList";
import { UserContext } from "../context/UserContext";
import "../styles/Profile.scss";

const Profile: React.FC = () => {
    const { user, loading } = useContext(UserContext);
    if (loading) return <div>Loading profile...</div>;
    if (!user) {
        return (
            <div className="must-login">
                You must be logged in to view this page.
            </div>
        );
    }
    const query = `?q=${user.username}&status=finished`;

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
                </div>

                <div className="vertical"></div>
                <div className="profile-gamelist">
                    <h1>Past games</h1>
                    <GameList renderFilters={false} query={query} />
                </div>
            </div>
        </div>
    );
};

export default Profile;
