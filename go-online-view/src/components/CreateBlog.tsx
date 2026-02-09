import React, { useEffect, useState } from "react";

type Props = {
    open: boolean;
    onClose: () => void;
    onCreated: (blogId: number | string) => void;
    onError?: (message: string) => void;
};

const CreateBlog: React.FC<Props> = ({ open, onClose, onCreated, onError }) => {
    const [id, setId] = useState<number | "">("");
    const [authorId, setAuthorId] = useState<number | "">("");
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");
    const [creating, setCreating] = useState(false);

    useEffect(() => {
        if (!open) return;
        setId("");
        setAuthorId("");
        setTitle("");
        setContent("");
    }, [open]);

    useEffect(() => {
        if (!open) return;
        const onKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape" && !creating) onClose();
        };
        window.addEventListener("keydown", onKeyDown);
        return () => window.removeEventListener("keydown", onKeyDown);
    }, [open, creating, onClose]);

    const createNewBlog = async () => {
        if (!id || !authorId || !title || !content) {
            console.log(id + " " + authorId + " " + title + " " + content);
            onError?.("All fields are required");
            return;
        }

        try {
            setCreating(true);

            const payload = {
                id,
                authorId,
                title,
                content,
            };

            const res = await fetch("/api/blogs", {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (!res.ok) {
                const text = await res.text();
                onError?.(text || "Failed to create blog");
                throw new Error(text);
            }

            const data = await res.json();
            onClose();
            onCreated(data.id);
        } catch (err) {
            console.error("Failed to create blog", err);
        } finally {
            setCreating(false);
        }
    };

    if (!open) return null;

    return (
        <div
            className="modal-backdrop"
            onMouseDown={() => !creating && onClose()}
            role="presentation"
        >
            <div
                className="modal"
                role="dialog"
                aria-modal="true"
                aria-label="Create blog"
                onMouseDown={(e) => e.stopPropagation()}
            >
                <div className="modal-title">Create Blog</div>

                <div className="modal-row">
                    <label>Blog ID</label>
                    <input
                        type="number"
                        value={id}
                        onChange={(e) =>
                            setId(
                                e.target.value === ""
                                    ? ""
                                    : Number(e.target.value),
                            )
                        }
                        disabled={creating}
                    />
                </div>

                <div className="modal-row">
                    <label>Author ID</label>
                    <input
                        type="number"
                        value={authorId}
                        onChange={(e) =>
                            setAuthorId(
                                e.target.value === ""
                                    ? ""
                                    : Number(e.target.value),
                            )
                        }
                        disabled={creating}
                    />
                </div>

                <div className="modal-row">
                    <label>Title</label>
                    <input
                        type="text"
                        value={title}
                        onChange={(e) => {
                            console.log(e.target.value);
                            setTitle(e.target.value);
                        }}
                        disabled={creating}
                    />
                </div>

                <div className="modal-row">
                    <label>Content</label>
                    <textarea
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        disabled={creating}
                        rows={6}
                        style={{ resize: "vertical" }}
                    />
                </div>

                <div className="modal-actions">
                    <button type="button" onClick={onClose} disabled={creating}>
                        Cancel
                    </button>

                    <button
                        type="button"
                        onClick={createNewBlog}
                        disabled={creating}
                    >
                        {creating ? "Creating…" : "Create"}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default CreateBlog;
