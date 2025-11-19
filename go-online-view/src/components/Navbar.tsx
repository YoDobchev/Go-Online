import { useContext, useEffect } from "react";
import { UserContext } from "../context/UserContext";

const Navbar: React.FC = () => {
    const { user, loading } = useContext(UserContext);

    useEffect(() => {
        console.log(user);
    }, [user]);
    return (
        <div className="navbar">
            {loading ? (
                "Loading..."
            ) : user ? (
                user.username
            ) : (
                <a href="/login">Login</a>
            )}
        </div>
    );
};

export default Navbar;
