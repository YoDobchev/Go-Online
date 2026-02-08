import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import "../styles/BlogContent.scss";
import BottomCenMsg from "../components/BottomCenMsg";
import CreateReply from "../components/CreateReply";

type Blog = {
  id: number;
  author_id: number;
  author_name: string;
  title: string;
  blog_content: string;
  published_at: string;
  updated_at?: string | null;
};

type BlogReply = {
  id: number;
  blog_id: number;
  author_id: number;
  author_name: string;
  reply_content: string;
  created_at: string;
};

export default function BlogContent() {
  const { id } = useParams<{ id: string }>();

  const [blog, setBlog] = useState<Blog | null>(null);
  const [replies, setReplies] = useState<BlogReply[]>([]);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [createOpen, setCreateOpen] = useState(false);
  const [show, setShow] = useState(false);
  const [msg, setMsg] = useState("");

  const triggerErrMsg = (text: string) => {
    setMsg(text);
    setShow(true);
  };

  const deleteReply = async (replyId: number) => {
    if (!id) return;

    try {
      const res = await fetch(`/api/blogs/${id}/replies`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify({ id: replyId }),
      });

      if (!res.ok) {
        throw new Error("Failed to delete reply");
      }

      setReplies((prev) => prev.filter((r) => r.id !== replyId));
    } catch (err) {
      console.error(err);
      triggerErrMsg("Could not delete reply");
    }
  };


  useEffect(() => {
    if (!id) return;

    Promise.all([
      fetch(`/api/blogs/${id}`).then((res) => {
        if (!res.ok) throw new Error("Failed to load blog");
        return res.json();
      }),
      fetch(`/api/blogs/${id}/replies`).then((res) => {
        if (!res.ok) throw new Error("Failed to load replies");
        return res.json();
      }),
    ])
      .then(([blogData, repliesData]) => {
        setBlog(blogData);
        setReplies(Array.isArray(repliesData) ? repliesData : []);
      })

      .catch((err) => setError(err.message))
      .finally(() => setLoading(false));
  }, [id]);

  const formatDate = (date?: string | null) =>
    date ? new Date(date).toISOString().slice(0, 10) : "-";

  if (loading) return <p>Loading blog…</p>;
  if (error) return <p style={{ color: "red" }}>Error: {error}</p>;
  if (!blog) return <p>Blog not found</p>;

  return (
    <div className="blogcontent">
      <a href="/blogs" className="back-link">Back to Blogs</a>

      <h1>{blog.title}</h1>

      <div className="meta">
        <div>Author: {blog.author_name} (ID: {blog.author_id})</div>
        <div>Published: {formatDate(blog.published_at)}</div>
        <div>Updated: {formatDate(blog.updated_at)}</div>
      </div>

      <div className="content">
        {blog.blog_content}
      </div>

      <div className="replies-header">
        <h2 className="replies-title">
          Replies ({replies.length})
        </h2>

        <button
          className="reply-add"
          onClick={() => setCreateOpen(true)}
          title="Add reply"
          aria-label="Add reply"
        >
          +
        </button>
      </div>


      {replies.length === 0 && (
        <div className="no-replies">No replies yet.</div>
      )}

      {replies.length > 0 && (
        <div className="replies">
          {replies.map((r) => (
            <div key={r.id} className="reply">
              <div className="reply-header">
                <div className="reply-meta">
                  {r.author_name} • {formatDate(r.created_at)}
                </div>

                <button
                  className="reply-delete"
                  onClick={() => deleteReply(r.id)}
                  title="Delete reply"
                  aria-label="Delete reply"
                >
                  ×
                </button>
              </div>

              <div className="reply-content">
                {r.reply_content}
              </div>
            </div>
          ))}
        </div>
      )}

      <CreateReply
        open={createOpen}
        blogId={Number(id)}
        onClose={() => setCreateOpen(false)}
        onCreated={() => {
          fetch(`/api/blogs/${id}/replies`)
            .then((res) => res.json())
            .then(setReplies);
        }}

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