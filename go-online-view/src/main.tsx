import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Route, Routes } from "react-router-dom";

import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Logout from "./pages/Logout";
import Game from "./pages/Game";
import Profile from "./pages/Profile";
import Blogs from "./pages/Blogs";
import Leaderboard from "./pages/Leaderboard";

import "./styles/index.scss";

import { AppProviders } from "./AppProviders";
import BlogContent from "./pages/BlogContent";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <AppProviders>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/logout" element={<Logout />} />
                <Route path="/game/:gameID" element={<Game />} />
                <Route path="/profile" element={<Profile />} />
                <Route path="/blogs" element={<Blogs />} />
                <Route path="/blogs/:id" element={<BlogContent />} />
                <Route path="/leaderboard" element={<Leaderboard />} />
            </Routes>
        </AppProviders>
    </StrictMode>
);
