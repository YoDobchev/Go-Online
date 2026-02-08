import React from "react";
import Navbar from "../components/Navbar";
import GameList from "../components/GameList";
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
};

const User: React.FC<ProfileProps> = ({ user, loading = false }) => {
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

export default User;
