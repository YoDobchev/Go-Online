import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { UserContext } from "../context/UserContext";
import "../styles/Dashboard.scss";
import Navbar from "../components/Navbar";

export default function Dashboard() {
    const { user } = useContext(UserContext);
    const navigate = useNavigate();

    const isModerator =
        !!user && (user.role === "admin" || user.role === "moderator");

    if (!isModerator) {
        return (
            <>
                <Navbar />
                <div className="reports-denied">
                    <h1 className="reports-denied__title">Access Denied</h1>
                    <p className="reports-denied__text">
                        You do not have permission to view this page.
                    </p>
                </div>
            </>
        );
    }

    return (
        <>
            <Navbar />
            <div className="dashboard">
                <h1>Admin Dashboard</h1>

                <div className="dashboard-actions">
                    <button onClick={() => navigate("/dashboard/users")}>
                        Users
                    </button>

                    <button onClick={() => navigate("/dashboard/reports")}>
                        Reports
                    </button>
                </div>
            </div>
        </>
    );
}
