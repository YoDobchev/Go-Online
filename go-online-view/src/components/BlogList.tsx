import { useEffect, useMemo, useState } from "react";
import "../styles/blogList.scss";

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

    // filters
    const [titleFilter, setTitleFilter] = useState("");
    const [authorFilter, setAuthorFilter] = useState("");
    const [fromDate, setFromDate] = useState("");
    const [toDate, setToDate] = useState("");

    useEffect(() => {
        fetch("/api/blogs")
            .then((res) => {
                if (!res.ok) throw new Error("Failed to load blogs");
                return res.json();
            })
            .then(setBlogs)
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

    const authors = Array.from(
        new Set(blogs.map((b) => b.author_name))
    );

    return (
        <div className="bloglist">
            <div className="bloglist-bar">
                <div className="title">Blogs</div>

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

                    {/* + button kept exactly as requested */}
                    <button className="plus">
                        <span>+</span>
                    </button>
                </div>
            </div>

            {/* Table */}
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
                            <div className="cell muted">{blog.id}</div>
                            <div className="cell">{blog.title}</div>
                            <div className="cell">{blog.author_name}</div>
                            <div className="cell">
                                {new Date(blog.published_at).toISOString().slice(0, 10)}
                            </div>
                            <div className="cell muted">
                                {blog.updated_at
                                    ? new Date(blog.updated_at).toISOString().slice(0, 10)
                                    : "—"}
                            </div>
                        </div>
                    ))}
            </div>
        </div>
    );
}
