import { useContext } from "react";
import Navbar from "../components/Navbar";
import GameList from "../components/GameList";
import { UserContext } from "../context/UserContext";

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
            <div>
                <h1>Profile</h1>
                <p>Email: {user.email}</p>
                <p>Username: {user.username}</p>
                <p>Elo: {user.elo}</p>
                <p>Rank: {user.rank}</p>
                <p>Role: {user.role}</p>

                <GameList  renderFilters={false} query={query}/>
            </div>
        </div>
    );
};

export default Profile;
