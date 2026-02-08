import { useEffect, useState } from "react";
import { useContext } from "react";
import { UserContext } from "../context/UserContext";


type Props = {
    open: boolean;
    blogId: number;
    onClose: () => void;
    onCreated: (reply: any) => void;
    onError?: (msg: string) => void;
};

export default function CreateReply({
    open,
    blogId,
    onClose,
    onCreated,
    onError,
}: Props) {
    const { user } = useContext(UserContext);
    const [content, setContent] = useState("");
    const [creating, setCreating] = useState(false);

    useEffect(() => {
        if (!open) return;
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

    const createReply = async () => {
        if (!user) {
            onError?.("You must be logged in to reply");
            return;
        }

        if (!content) {
            onError?.("Content is required");
            return;
        }

        try {
            setCreating(true);

            const res = await fetch(`/api/blogs/${blogId}/replies`, {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    blogId,
                    replyContent: content,
                }),
            });

            if (!res.ok) {
                const text = await res.text();
                throw new Error(text || "Failed to create reply");
            }

            const data = await res.json();
            onCreated(data);
            onClose();
        } catch (err: any) {
            console.error(err);
            onError?.(err.message || "Failed to create reply");
        } finally {
            setCreating(false);
        }
    };

    if (!open) return null;

    return (
        <div className="modal-backdrop" onMouseDown={onClose}>
            <div
                className="modal"
                onMouseDown={(e) => e.stopPropagation()}
                role="dialog"
                aria-modal="true"
            >
                <div className="modal-title">Add Reply</div>

                <div className="modal-row">
                    <label>Reply</label>
                    <textarea
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        rows={5}
                        disabled={creating}
                    />
                </div>

                <div className="modal-actions">
                    <button onClick={onClose} disabled={creating}>
                        Cancel
                    </button>
                    <button onClick={createReply} disabled={creating}>
                        {creating ? "Posting…" : "Post"}
                    </button>
                </div>
            </div>
        </div>
    );
}