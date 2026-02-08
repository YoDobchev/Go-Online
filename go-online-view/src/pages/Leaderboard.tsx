import { useContext } from "react";
import LeaderboardList from "../components/LeaderboardList";
import Navbar from "../components/Navbar"
import { UserContext } from "../context/UserContext";

const Leaderboard: React.FC = () => {
    const { user, loading } = useContext(UserContext);
    if (loading) return <div>Loading profile...</div>;
    if (!user) {
        return (
            <div className="must-login">
                <Navbar />
                You must be logged in to view this page.
            </div>
        );
    }
    return (
        <>
            <Navbar />
            <div className="leaderboard-box">
                <h1>Leaderboard</h1>
                <LeaderboardList />
            </div>
        </>
    )
};

export default Leaderboard;