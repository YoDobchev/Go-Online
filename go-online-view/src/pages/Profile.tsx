import { useContext } from "react";
import User from "../components/User";
import { UserContext } from "../context/UserContext";

export default function ProfilePage() {
    const { user, loading } = useContext(UserContext);
    return <User user={user} loading={loading} />;
}
