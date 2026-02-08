import { useContext, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { UserContext } from "../context/UserContext";
import "../styles/Dashboard.scss";

export default function Dashboard() {
  const { user } = useContext(UserContext);
  const navigate = useNavigate();

  if (!user) return null;

  return (
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
  );
}