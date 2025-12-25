import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Route, Routes } from "react-router-dom";

import Home from "./pages/Home";
import Login from "./pages/Login";
import Register from "./pages/Register";
import Logout from "./pages/Logout";

import "./styles/index.scss";

import { AppProviders } from "./AppProviders";

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <AppProviders>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/logout" element={<Logout />} />
            </Routes>
        </AppProviders>
    </StrictMode>
);
