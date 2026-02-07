import React, { useEffect, useState } from "react";
import { API_BASE } from "../config";

type GameStatus = { id: number | string };

type CreateSize = "9" | "13" | "19";
type PlayAs = 0 | 1 | 2;

function parseCreateSize(v: string): CreateSize {
    return v === "9" || v === "13" || v === "19" ? v : "19";
}
function parsePlayAs(v: string): PlayAs {
    return v === "0" || v === "1" || v === "2" ? (Number(v) as PlayAs) : 0;
}

type Props = {
    open: boolean;
    onClose: () => void;
    onCreated: (gameId: number | string) => void;
    onError?: (message: string) => void;

    defaultPlayAs?: PlayAs;
    defaultSize?: CreateSize;
    defaultRanked?: boolean;
};

const CreateGame: React.FC<Props> = ({
    open,
    onClose,
    onCreated,
    onError,
    defaultPlayAs = 0,
    defaultSize = "19",
    defaultRanked = true,
}) => {
    const [playAs, setPlayAs] = useState<PlayAs>(defaultPlayAs);
    const [createSize, setCreateSize] = useState<CreateSize>(defaultSize);
    const [createRanked, setCreateRanked] = useState<boolean>(defaultRanked);
    const [vsAI, setVsAI] = useState(false);
    const [creating, setCreating] = useState(false);

    useEffect(() => {
        if (!open) return;
        setPlayAs(defaultPlayAs);
        setCreateSize(defaultSize);
        setCreateRanked(defaultRanked);
    }, [open, defaultPlayAs, defaultSize, defaultRanked]);

    useEffect(() => {
        if (!open) return;

        const onKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape" && !creating) onClose();
        };

        window.addEventListener("keydown", onKeyDown);
        return () => window.removeEventListener("keydown", onKeyDown);
    }, [open, creating, onClose]);

    const createNewGame = async () => {
        try {
            setCreating(true);

            const payload = {
                playAs,
                boardSize: Number(createSize),
                ranked: createRanked,
                vsAI
            };

            const res = await fetch(`${API_BASE}/game/`, {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (!res.ok) {
                const errorText = await res.text();

                if (errorText.includes("user already in a game")) {
                    onError?.("You are already in a game!");
                }

                throw new Error(errorText || "Failed to create new game");
            }

            const data: GameStatus = await res.json();

            onClose();
            onCreated(data.id);
        } catch (err) {
            console.error("Failed to create new game", err);
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
                aria-label="Create game"
                onMouseDown={(e) => e.stopPropagation()}
            >
                <div className="modal-title">Create game</div>

                <div className="modal-row">
                    <label>Play as</label>
                    <select
                        value={playAs}
                        onChange={(e) => setPlayAs(parsePlayAs(e.target.value))}
                        disabled={creating}
                    >
                        <option value="0">Black</option>
                        <option value="1">White</option>
                        <option value="2">Random</option>
                    </select>
                </div>

                <div className="modal-row">
                    <label>Board size</label>
                    <select
                        value={createSize}
                        onChange={(e) =>
                            setCreateSize(parseCreateSize(e.target.value))
                        }
                        disabled={creating}
                    >
                        <option value="9">9x9</option>
                        <option value="13">13x13</option>
                        <option value="19">19x19</option>
                    </select>
                </div>

                <div className="modal-row">
                    <label>Opponent</label>
                    <select
                        value={vsAI ? "ai" : "human"}
                        onChange={(e) => setVsAI(e.target.value === "ai")}
                        disabled={creating}
                    >
                        <option value="human">Human</option>
                        <option value="ai">KataGo (AI)</option>
                    </select>
                </div>

                <div className="modal-row">
                    <label>Type</label>
                    <select
                        value={createRanked ? "ranked" : "unranked"}
                        onChange={(e) =>
                            setCreateRanked(e.target.value === "ranked")
                        }
                        disabled={creating}
                    >
                        <option value="ranked">Ranked</option>
                        <option value="unranked">Unranked</option>
                    </select>
                </div>

                <div className="modal-actions">
                    <button type="button" onClick={onClose} disabled={creating}>
                        Cancel
                    </button>

                    <button
                        type="button"
                        onClick={createNewGame}
                        disabled={creating}
                    >
                        {creating ? "Creating…" : "Create"}
                    </button>
                </div>
            </div>
        </div>
    );
};

export default CreateGame;
