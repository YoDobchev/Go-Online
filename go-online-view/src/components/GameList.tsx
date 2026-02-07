import React, { useEffect, useMemo, useState } from "react";
import { API_BASE } from "../config";
import "../styles/gamelist.scss";

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

                const res = await fetch(`${API_BASE}/search/games?${queryString}`, {
                    credentials: "include",
                    signal: controller.signal,
                });

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

    return (
        <div className="gamelist">
            <div className="gamelist-bar">
                {/* <div className="bar-left">
                    <div className="title">Games</div>
                </div> */}

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
                        <option value="casual">Casual</option>
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

                    <a
                        className="plus"
                        href="/games/create"
                        aria-label="Create game"
                        title="Create game"
                    >
                        <span>+</span>
                    </a>
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
                            href={`/games/${g.id}`}
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
                                      : "Casual"}
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
        </div>
    );
};

export default GameList;
