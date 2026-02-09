import { useContext, useEffect, useMemo, useState } from "react";
import "../styles/BlogList.scss";
import { useNavigate } from "react-router-dom";
import BottomCenMsg from "./BottomCenMsg";
import CreateBlog from "./CreateBlog";
import { UserContext } from "../context/UserContext";

type Blog = {
    id: number;
    title: string;
    author_name: string;
    published_at: string;
    updated_at?: string | null;
};

export default function BlogsList() {
    const [blogs, setBlogs] = useState<Blog[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [titleFilter, setTitleFilter] = useState("");
    const [authorFilter, setAuthorFilter] = useState("");
    const [fromDate, setFromDate] = useState("");
    const [toDate, setToDate] = useState("");

    const navigate = useNavigate();

    const [createOpen, setCreateOpen] = useState(false);
    const [show, setShow] = useState(false);
    const [msg, setMsg] = useState("");

    const { user } = useContext(UserContext);

    const triggerErrMsg = (text: string) => {
        setMsg(text);
        setShow(true);
    };

    useEffect(() => {
        fetch("/api/blogs")
            .then((res) => {
                if (!res.ok) throw new Error("Failed to load blogs");
                return res.json();
            })
            .then((data) => {
                const list = Array.isArray(data)
                    ? data
                    : Array.isArray(data?.blogs)
                      ? data.blogs
                      : [];
                setBlogs(list);
            })
            .catch((err) => setError(err.message))
            .finally(() => setLoading(false));
    }, []);

    const filteredBlogs = useMemo(() => {
        return blogs.filter((blog) => {
            const titleMatch = blog.title
                .toLowerCase()
                .includes(titleFilter.toLowerCase());

            const authorMatch = authorFilter
                ? blog.author_name === authorFilter
                : true;

            const published = new Date(blog.published_at).getTime();

            const fromMatch = fromDate
                ? published >= new Date(fromDate).getTime()
                : true;

            const toMatch = toDate
                ? published <= new Date(toDate).getTime()
                : true;

            return titleMatch && authorMatch && fromMatch && toMatch;
        });
    }, [blogs, titleFilter, authorFilter, fromDate, toDate]);

    const authors = Array.from(new Set(blogs.map((b) => b.author_name)));

    return (
        <div className="bloglist">
            <div className="bloglist-bar">
                <div className="title">Filter</div>

                <div className="bar-right">
                    <input
                        placeholder="Title keywords…"
                        value={titleFilter}
                        onChange={(e) => setTitleFilter(e.target.value)}
                    />

                    <select
                        value={authorFilter}
                        onChange={(e) => setAuthorFilter(e.target.value)}
                    >
                        <option value="">All authors</option>
                        {authors.map((a) => (
                            <option key={a} value={a}>
                                {a}
                            </option>
                        ))}
                    </select>

                    <input
                        type="date"
                        value={fromDate}
                        onChange={(e) => setFromDate(e.target.value)}
                    />

                    <input
                        type="date"
                        value={toDate}
                        onChange={(e) => setToDate(e.target.value)}
                    />

                    {user && user.role === "admin" && (
                        <button
                            className="plus"
                            onClick={() => setCreateOpen(true)}
                            aria-label="Create blog"
                            title="Create blog"
                        >
                            <span>+</span>
                        </button>
                    )}
                </div>
            </div>

            <div className="bloglist-box">
                <div className="bloglist-header">
                    <div>ID</div>
                    <div>Title</div>
                    <div>Author</div>
                    <div>Published</div>
                    <div>Updated</div>
                </div>

                {loading && <div className="state">Loading blogs…</div>}
                {error && <div className="state error">{error}</div>}
                {!loading && !error && filteredBlogs.length === 0 && (
                    <div className="state">No blogs match the filters</div>
                )}

                {!loading &&
                    !error &&
                    filteredBlogs.map((blog) => (
                        <div key={blog.id} className="blogrow">
                            <a
                                href={`/blogs/${blog.id}`}
                                style={{
                                    display: "contents",
                                    color: "inherit",
                                    textDecoration: "none",
                                }}
                            >
                                <div className="cell muted">{blog.id}</div>
                                <div className="cell">{blog.title}</div>
                                <div className="cell">{blog.author_name}</div>
                                <div className="cell">
                                    {new Date(blog.published_at)
                                        .toISOString()
                                        .slice(0, 10)}
                                </div>
                                <div className="cell muted">
                                    {blog.updated_at
                                        ? new Date(blog.updated_at)
                                              .toISOString()
                                              .slice(0, 10)
                                        : "-"}
                                </div>
                            </a>
                        </div>
                    ))}
            </div>

            <CreateBlog
                open={createOpen}
                onClose={() => setCreateOpen(false)}
                onCreated={(id) => navigate(`/blogs/${id}`)}
                onError={triggerErrMsg}
            />

            <BottomCenMsg
                visible={show}
                message={msg}
                backgroundColor="#e90e0e"
                textColor="#f7f7f7"
                timeAfterFadeMs={2000}
                onClose={() => setShow(false)}
            />
        </div>
    );
}
