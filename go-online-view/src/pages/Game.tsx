import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { API_BASE } from "../config";

type Stone = 0 | 1 | 2;

function GoBoardSVG({
    board,
    onPlay,
    interactive = false,
}: {
    board: Board;
    onPlay?: (row: number, col: number) => void;
    interactive?: boolean;
}) {
    const N = board.size;
    const padding = 24;
    const cell = 28;
    const w = padding * 2 + cell * (N - 1);
    const h = padding * 2 + cell * (N - 1);

    const x = (col: number) => padding + col * cell;
    const y = (row: number) => padding + row * cell;

    const starCoords = N === 19 ? [3, 9, 15] : [];
    const starPoints =
        starCoords.length > 0
            ? starCoords.flatMap((r) => starCoords.map((c) => ({ r, c })))
            : [];

    const handleClick = (row: number, col: number) => {
        if (!interactive) return;
        if (!onPlay) return;
        onPlay(row, col);
    };

    return (
        <svg
            width={w}
            height={h}
            viewBox={`0 0 ${w} ${h}`}
            style={{
                display: "block",
                maxWidth: "100%",
                height: "auto",
                userSelect: "none",
            }}
        >
            <rect x={0} y={0} width={w} height={h} fill="#deb887" />

            {Array.from({ length: N }).map((_, i) => (
                <g key={i}>
                    <line
                        x1={x(i)}
                        y1={y(0)}
                        x2={x(i)}
                        y2={y(N - 1)}
                        stroke="rgba(0,0,0,0.75)"
                        strokeWidth={1}
                    />
                    <line
                        x1={x(0)}
                        y1={y(i)}
                        x2={x(N - 1)}
                        y2={y(i)}
                        stroke="rgba(0,0,0,0.75)"
                        strokeWidth={1}
                    />
                </g>
            ))}

            {starPoints.map((p, idx) => (
                <circle
                    key={idx}
                    cx={x(p.c)}
                    cy={y(p.r)}
                    r={3}
                    fill="rgba(0,0,0,0.85)"
                />
            ))}

            {Array.from({ length: N }).map((_, r) =>
                Array.from({ length: N }).map((__, c) => (
                    <circle
                        key={`hit-${r}-${c}`}
                        cx={x(c)}
                        cy={y(r)}
                        r={cell * 0.45}
                        fill="transparent"
                        style={{ cursor: interactive ? "pointer" : "default" }}
                        onClick={() => handleClick(r, c)}
                    />
                )),
            )}

            {board.squares.map((row, r) =>
                row.map((v, c) => {
                    const stone = v as Stone;
                    if (stone === 0) return null;

                    const isBlack = stone === 1;
                    return (
                        <g key={`stone-${r}-${c}`}>
                            <circle
                                cx={x(c) + 1.2}
                                cy={y(r) + 1.2}
                                r={cell * 0.42}
                                fill="rgba(0,0,0,0.25)"
                            />
                            <circle
                                cx={x(c)}
                                cy={y(r)}
                                r={cell * 0.42}
                                fill={isBlack ? "#f5f5f5" : "#111"}
                                stroke={isBlack ? "#bbb" : "#000"}
                                strokeWidth={1}
                            />
                            <circle
                                cx={x(c) - cell * 0.12}
                                cy={y(r) - cell * 0.12}
                                r={cell * 0.12}
                                fill={
                                    isBlack
                                        ? "rgba(255,255,255,0.55)"
                                        : "rgba(255,255,255,0.12)"
                                }
                            />
                        </g>
                    );
                }),
            )}
        </svg>
    );
}

interface Board {
    size: number;
    squares: number[][];
}
interface GameInfo {
    players: [string, string];
    turn: number;
    board: Board;
}

type ServerMsg =
    | { type: "hello"; data: { role: "player" | "spectator"; user: string } }
    | { type: "game_snapshot"; data: GameInfo }
    | { type: "move_played"; data: { by: string; row: number; col: number } }
    | { type: "error"; data: string };

const Game: React.FC = () => {
    const { gameID } = useParams<{ gameID: string }>();
    const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
    const [role, setRole] = useState<"player" | "spectator">("spectator");
    const wsRef = useRef<WebSocket | null>(null);

    useEffect(() => {
        if (!gameID) return;

        const wsBase = API_BASE.replace(/^http/, "ws");
        const wsUrl = `${wsBase}/game/${gameID}/ws`;

        const ws = new WebSocket(wsUrl);
        wsRef.current = ws;

        ws.onopen = () => {
            console.log("WS connected");
        };

        ws.onmessage = (ev) => {
            const msg = JSON.parse(ev.data) as ServerMsg;

            if (msg.type === "hello") {
                setRole(msg.data.role);
            } else if (msg.type === "game_snapshot") {
                setGameInfo(msg.data);
            } else if (msg.type === "error") {
                console.error("Server error:", msg.data);
            }
        };

        ws.onerror = (e) => {
            console.error("WS error", e);
        };

        ws.onclose = () => {
            console.log("WS closed");
        };

        return () => {
            ws.close();
        };
    }, [gameID]);

    const sendMove = (row: number, col: number) => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;

        ws.send(
            JSON.stringify({
                type: "play_move",
                data: { row, col },
            }),
        );
    };

    if (!gameInfo) return <div>Loading...</div>;

    return (
        <div>
            <div>Role: {role}</div>
            {role === "spectator" && <div>You are spectating (read-only).</div>}

            <div>
                Players: {gameInfo.players[0]} vs{" "}
                {gameInfo.players[1] || "(waiting)"}
            </div>
            <div>Turn: {gameInfo.turn}</div>
            <GoBoardSVG
                board={gameInfo.board}
                interactive={role === "player"}
                onPlay={(r, c) => sendMove(r, c)}
            />
        </div>
    );
};

export default Game;
