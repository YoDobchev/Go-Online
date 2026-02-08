import Navbar from "../components/Navbar";
import BlogList from "../components/BlogList";

const Blogs: React.FC = () => {

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
