import React, { useEffect, useMemo, useState } from "react";
import { API_BASE } from "../config";
import "../styles/gamelist.scss";
import { useNavigate } from "react-router-dom";

type GameStatus = { id: number | string };

type Game = {
    id: number | string;
    name?: string;
    status?: "open" | "running" | "finished" | string;
    boardSize?: number;
    ranked?: boolean;
    players?: number;
    createdAt?: string;
};

type StatusFilter = "all" | "open" | "running" | "finished";
type RankedFilter = "all" | "ranked" | "unranked";
type SizeFilter = "all" | "9" | "13" | "19";

// type GameListResponse =
//   | Game[]
//   | {
//       games?: Game[];
//       items?: Game[];
//       data?: Game[];
//       page?: number;
//       totalPages?: number;
//       pages?: number;
//     };

function isObject(v: unknown): v is Record<string, unknown> {
    return typeof v === "object" && v !== null;
}

function toGameArray(data: unknown): Game[] {
    if (Array.isArray(data)) return data as Game[];

    if (isObject(data)) {
        const maybe = data.games ?? data.items ?? data.data;
        if (Array.isArray(maybe)) return maybe as Game[];
    }

    return [];
}

function toTotalPages(data: unknown): number {
    if (!isObject(data)) return 1;

    const tp = data.totalPages;
    const pages = data.pages;

    if (typeof tp === "number" && Number.isFinite(tp) && tp >= 1) return tp;
    if (typeof pages === "number" && Number.isFinite(pages) && pages >= 1)
        return pages;

    return 1;
}

function parseStatus(v: string): StatusFilter {
    return v === "open" || v === "running" || v === "finished" || v === "all"
        ? v
        : "all";
}
function parseRanked(v: string): RankedFilter {
    return v === "ranked" || v === "unranked" || v === "all" ? v : "all";
}
function parseSize(v: string): SizeFilter {
    return v === "9" || v === "13" || v === "19" || v === "all" ? v : "all";
}

const GameList: React.FC = () => {
    const [games, setGames] = useState<Game[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);

    const [q, setQ] = useState("");
    const [status, setStatus] = useState<StatusFilter>("all");
    const [ranked, setRanked] = useState<RankedFilter>("all");
    const [size, setSize] = useState<SizeFilter>("all");

    const queryString = useMemo(() => {
        const params = new URLSearchParams();
        params.set("p", String(page));
        if (q.trim()) params.set("q", q.trim());
        if (status !== "all") params.set("status", status);
        if (ranked !== "all")
            params.set("ranked", ranked === "ranked" ? "true" : "false");
        if (size !== "all") params.set("size", size);
        return params.toString();
    }, [page, q, status, ranked, size]);

    useEffect(() => {
        const controller = new AbortController();

        (async () => {
            try {
                setLoading(true);
                setError(null);

                const res = await fetch(
                    `${API_BASE}/search/games?${queryString}`,
                    {
                        credentials: "include",
                        signal: controller.signal,
                    },
                );

                if (!res.ok)
                    throw new Error(`Failed to load games (${res.status})`);

                const data: unknown = await res.json();
                const items = toGameArray(data);
                const tp = toTotalPages(data);

                setGames(items);
                setTotalPages(tp);
            } catch (e: unknown) {
                if (e instanceof DOMException && e.name === "AbortError")
                    return;
                const msg =
                    e instanceof Error ? e.message : "Something went wrong";
                setError(msg);
            } finally {
                setLoading(false);
            }
        })();

        return () => controller.abort();
    }, [queryString]);

    const resetToPage1 = () => setPage(1);

    const navigate = useNavigate();

    const [createOpen, setCreateOpen] = useState(false);
    const [playAs, setPlayAs] = useState(0);
    const [createSize, setCreateSize] = useState<SizeFilter>("19");
    const [createRanked, setCreateRanked] = useState(true);
    const [creating, setCreating] = useState(false);

    useEffect(() => {
        if (!createOpen) return;
        const onKeyDown = (e: KeyboardEvent) => {
            if (e.key === "Escape") setCreateOpen(false);
        };
        window.addEventListener("keydown", onKeyDown);
        return () => window.removeEventListener("keydown", onKeyDown);
    }, [createOpen]);

    const createNewGame = async () => {
        try {
            setCreating(true);

            const payload = {
                playAs: playAs,
                boardSize: Number(createSize),
                ranked: createRanked,
            };

            const res = await fetch(`${API_BASE}/game/`, {
                method: "POST",
                credentials: "include",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(payload),
            });

            if (!res.ok) throw new Error("Failed to create new game");

            const data: GameStatus = await res.json();

            setCreateOpen(false);

            navigate(`/game/${data.id}`);
        } catch (err) {
            console.error("Failed to create new game", err);
        } finally {
            setCreating(false);
        }
    };

    return (
        <div className="gamelist">
            <div className="gamelist-bar">
                <div className="bar-right">
                    <input
                        value={q}
                        onChange={(e) => {
                            setQ(e.target.value);
                            resetToPage1();
                        }}
                        placeholder="Search…"
                        aria-label="Search games"
                    />

                    <select
                        value={status}
                        onChange={(e) => {
                            setStatus(parseStatus(e.target.value));
                            resetToPage1();
                        }}
                        aria-label="Filter by status"
                    >
                        <option value="all">Status</option>
                        <option value="open">Open</option>
                        <option value="running">Running</option>
                        <option value="finished">Finished</option>
                    </select>

                    <select
                        value={ranked}
                        onChange={(e) => {
                            setRanked(parseRanked(e.target.value));
                            resetToPage1();
                        }}
                        aria-label="Filter by type"
                    >
                        <option value="all">Type</option>
                        <option value="ranked">Ranked</option>
                        <option value="unranked">unranked</option>
                    </select>

                    <select
                        value={size}
                        onChange={(e) => {
                            setSize(parseSize(e.target.value));
                            resetToPage1();
                        }}
                        aria-label="Filter by board size"
                    >
                        <option value="all">Size</option>
                        <option value="9">9x9</option>
                        <option value="13">13x13</option>
                        <option value="19">19x19</option>
                    </select>

                    <button
                        type="button"
                        className="plus"
                        onClick={() => setCreateOpen(true)}
                        aria-label="Create game"
                        title="Create game"
                    >
                        <span>+</span>
                    </button>
                </div>
            </div>

            <div className="gamelist-box">
                <div className="gamelist-header">
                    <div>ID</div>
                    <div>Name</div>
                    <div>Status</div>
                    <div>Size</div>
                    <div>Type</div>
                    <div>Players</div>
                </div>

                {loading && <div className="state">Loading…</div>}
                {error && <div className="state error">{error}</div>}
                {!loading && !error && games.length === 0 && (
                    <div className="state">No games found.</div>
                )}

                {!loading &&
                    !error &&
                    games.map((g) => (
                        <a
                            key={g.id}
                            className="gamerow"
                            href={`/game/${g.id}`}
                        >
                            <div className="cell muted">{g.id}</div>
                            <div className="cell">{g.name ?? "—"}</div>
                            <div className="cell">{g.status ?? "—"}</div>
                            <div className="cell">
                                {g.boardSize
                                    ? `${g.boardSize}x${g.boardSize}`
                                    : "—"}
                            </div>
                            <div className="cell">
                                {g.ranked === undefined
                                    ? "—"
                                    : g.ranked
                                      ? "Ranked"
                                      : "Unranked"}
                            </div>
                            <div className="cell">{g.players ?? "—"} / 2</div>
                        </a>
                    ))}

                <div className="pager">
                    <button
                        onClick={() => setPage((p) => Math.max(1, p - 1))}
                        disabled={page <= 1 || loading}
                    >
                        Prev
                    </button>

                    <span className="page">
                        {page} / {Math.max(1, totalPages)}
                    </span>

                    <button
                        onClick={() =>
                            setPage((p) => Math.min(totalPages, p + 1))
                        }
                        disabled={page >= totalPages || loading}
                    >
                        Next
                    </button>
                </div>
            </div>

            {createOpen && (
                <div
                    className="modal-backdrop"
                    onMouseDown={() => !creating && setCreateOpen(false)}
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
                                onChange={(e) =>
                                    setPlayAs(parseInt(e.target.value))
                                }
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
                                    setCreateSize(parseSize(e.target.value))
                                }
                                disabled={creating}
                            >
                                <option value="9">9x9</option>
                                <option value="13">13x13</option>
                                <option value="19">19x19</option>
                            </select>
                        </div>

                        <div className="modal-row">
                            <label>Type</label>
                            <select
                                value={createRanked ? "ranked" : "unrakned"}
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
                            <button
                                type="button"
                                onClick={() => setCreateOpen(false)}
                                disabled={creating}
                            >
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
            )}
        </div>
    );
};

export default GameList;
