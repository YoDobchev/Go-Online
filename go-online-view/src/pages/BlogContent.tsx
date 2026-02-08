import { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";

type Blog = {
  id: number;
  author_id: number;
  author_name: string;
  title: string;
  blog_content: string;
  published_at: string;
  updated_at?: string | null;
};

export default function BlogContent() {
  const { id } = useParams<{ id: string }>();
  const [blog, setBlog] = useState<Blog | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    fetch(`/api/blogs/${id}`)
      .then((res) => {
        if (!res.ok) throw new Error("Failed to load blog");
        return res.json();
      })
      .then(setBlog)
      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [id]);

  const formatDate = (date?: string | null) =>
    date ? new Date(date).toISOString().slice(0, 10) : "—";

  if (loading) return <p>Loading blog…</p>;
  if (error) return <p style={{ color: "red" }}>Error: {error}</p>;
  if (!blog) return <p>Blog not found</p>;

  return (
    <div className="blogcontent" style={{ maxWidth: "720px", margin: "2rem auto", padding: "1rem" }}>
      <a href="/blogs" style={{ display: "inline-block", marginBottom: "1rem" }}>Back to Blogs</a>

      <h1 style={{ marginBottom: "0.5rem" }}>{blog.title}</h1>
      <div style={{ marginBottom: "1rem", color: "#666" }}>
        <span>Author: {blog.author_name} (ID: {blog.author_id})</span><br />
        <span>Published: {formatDate(blog.published_at)}</span><br />
        <span>Updated: {formatDate(blog.updated_at)}</span>
      </div>

      <div style={{ lineHeight: "1.6", fontSize: "1rem", whiteSpace: "pre-wrap", border: "1px solid #eee", padding: "1rem", borderRadius: "10px", background: "#fafafa" }}>
        {blog.blog_content}
      </div>
    </div>
  );
}