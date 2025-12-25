import { useEffect, useContext } from "react";
import { useNavigate } from "react-router-dom";
import { API_BASE } from "../config";
import { UserContext } from "../context/UserContext";

export default function Logout() {
    const navigate = useNavigate();
    const { setUser } = useContext(UserContext);

    useEffect(() => {
        fetch(`${API_BASE}/auth/logout`, {
            method: "DELETE",
            credentials: "include",
        }).finally(() => {
            setUser(null);
            navigate("/", { replace: true });
        });
    }, [navigate, setUser]);

    return null;
}
