import { useContext } from "react";
import Navbar from "../components/Navbar";
import { UserContext } from "../context/UserContext";
import BlogList from "../components/BlogList";

const Blogs: React.FC = () => {
    const { user, loading } = useContext(UserContext);
    if (loading) return <div>Loading profile...</div>;
    if (!user) {
        return (
            <div className="must-login">
                You must be logged in to view this page.
            </div>
        );
    }

    return (
        <div>
            <Navbar />
            <div>
                <h1>Blogs</h1>
                <BlogList/>
            </div>
        </div>
    );
};

export default Blogs;
