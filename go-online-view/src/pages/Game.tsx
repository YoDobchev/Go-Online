import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { API_BASE } from "../config";
import GoBoardSVG from "../components/GoBoardSVG";
import type { Board } from "../components/GoBoardSVG";
import GameHistory from "../components/GameHistory";

interface GameInfo {
    players: [string, string];
    turn: number;
    board: Board;
    moveNum: number;
}

interface Status {
    role: "player" | "spectator";
    seat: "black" | "white" | "spectator";
}

interface ClockSnapshot {
    player: number; // 1 = black, 2 = white
    main_remaining_ms: number;
    byo_remaining_ms: number;
    byo_periods_left: number;
    server_time: number;
}

type ServerMsg =
    | {
          type: "hello";
          data: {
              role: "player" | "spectator";
              seat: "black" | "white" | "spectator";
          };
      }
    | { type: "sync"; data: GameInfo }
    | {
          type: "game_ended";
          data: { white_points: number; black_points: number; winner: number };
      }
    | { type: "error"; data: string }
    | { type: "clock_update"; data: ClockSnapshot }
    | { type: "timeout"; data: { loser: number } };

const Game: React.FC = () => {
    const { gameID } = useParams<{ gameID: string }>();
    const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
    const [status, setStatus] = useState<Status | null>(null);

    const [viewMoveNum, setViewMoveNum] = useState<number | null>(null);
    const [viewBoard, setViewBoard] = useState<Board | null>(null);

    const [clocks, setClocks] = useState<Record<number, ClockSnapshot>>({});
    const [clockReceivedAt, setClockReceivedAt] = useState<number>(Date.now());

    const [pendingMove, setPendingMove] = useState(false);

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
                setStatus(msg.data);
            } else if (msg.type === "sync") {
                setGameInfo(msg.data);
                setPendingMove(false);
                setClocks((prev) => ({
                    1: prev[1] ?? {
                        player: 1,
                        main_remaining_ms: 0,
                        byo_remaining_ms: 0,
                        byo_periods_left: 0,
                        server_time: Date.now(),
                    },
                    2: prev[2] ?? {
                        player: 2,
                        main_remaining_ms: 0,
                        byo_remaining_ms: 0,
                        byo_periods_left: 0,
                        server_time: Date.now(),
                    },
                }));
            } else if (msg.type === "game_ended") {
                console.log(msg.type, msg.data);
                console.log("gameended");
            } else if (msg.type === "clock_update") {
                setClocks((prev) => ({
                    ...prev,
                    [msg.data.player]: msg.data,
                }));
                setClockReceivedAt(Date.now());
            } else if (msg.type === "timeout") {
                console.log(`${msg.data.loser} ran out of time`);
            } else if (msg.type === "error") {
                setPendingMove(false);
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

    useEffect(() => {
        const id = setInterval(() => {
            setClocks((c) => ({ ...c }));
        }, 250);
        return () => clearInterval(id);
    }, []);

    const sendMove = (row: number, col: number) => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;
        if (!canPlay) return;

        setPendingMove(true);
        ws.send(JSON.stringify({ type: "play_move", data: { row, col } }));
    };

    const sendPass = () => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;
        if (!canPlay) return;

        setPendingMove(true);
        ws.send(
            JSON.stringify({ type: "play_move", data: { row: -1, col: -1 } }),
        );
    };

    const fetchBoardAtMove = async (moveNum: number): Promise<Board> => {
        const res = await fetch(
            `${API_BASE}/game/${gameID}/state?moveNum=${moveNum}`,
        );
        if (!res.ok) throw new Error(await res.text());

        const snap = (await res.json()) as Board;
        return snap;
    };

    const handleSelectMove = async (moveNum: number) => {
        setViewMoveNum(moveNum);
        try {
            const b = await fetchBoardAtMove(moveNum);
            setViewBoard(b);
        } catch (e) {
            console.error("Failed to load snapshot:", e);
        }
    };

    const getDisplayedTime = (player: number) => {
        const clock = clocks[player];
        if (!clock) return null;

        let remaining =
            clock.main_remaining_ms > 0
                ? clock.main_remaining_ms
                : clock.byo_remaining_ms;

        if (gameInfo?.turn === player) {
            remaining -= Date.now() - clockReceivedAt;
        }

        return Math.max(0, remaining);
    };

    const formatTime = (ms: number) => {
        const totalSeconds = Math.ceil(ms / 1000);
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${minutes}:${seconds.toString().padStart(2, "0")}`;
    };

    const backToLive = () => {
        setViewMoveNum(null);
        setViewBoard(null);
    };

    if (!gameInfo || !status) return <div>Loading...</div>;

    const maxMoveNum = gameInfo.moveNum ?? 0;
    const boardToShow = viewBoard ?? gameInfo.board;
    const isLive = viewMoveNum === null;

    const opponentJoined = Boolean(gameInfo.players[1]);
    const isLobbyWaiting = status.role === "player" && !opponentJoined;

    const myTurn =
        status.seat === "black" ? 2 : status.seat === "white" ? 1 : 0;

    const canPlay =
        status.role === "player" &&
        isLive &&
        opponentJoined &&
        !pendingMove &&
        gameInfo.turn === myTurn;

    const PlayerClock: React.FC<{
        label: string;
        player: number;
        active: boolean;
    }> = ({ label, player, active }) => {
        const remaining = getDisplayedTime(player);
        const clock = clocks[player];

        if (!clock) return null;

        return (
            <div
                style={{
                    padding: "6px 10px",
                    borderRadius: 6,
                    border: "1px solid #ddd",
                    background: active ? "#eef6ff" : "#fafafa",
                    marginBottom: 6,
                }}
            >
                <div style={{ fontWeight: 600 }}>{label}</div>
                <div style={{ fontSize: 18 }}>
                    {remaining !== null ? formatTime(remaining) : "--:--"}
                </div>

                {clock.main_remaining_ms <= 0 && (
                    <div style={{ fontSize: 12, color: "#555" }}>
                        Byo-yomi: {clock.byo_periods_left} ×{" "}
                        {formatTime(clock.byo_remaining_ms)}
                    </div>
                )}
            </div>
        );
    };

    return (
        <div style={{ display: "flex", gap: 16, alignItems: "flex-start" }}>
            <div>
                <div>Role: {status.role}</div>

                {status.role === "spectator" && (
                    <div>You are spectating (read-only).</div>
                )}

                {isLobbyWaiting && (
                    <div
                        style={{
                            margin: "8px 0",
                            padding: "8px 10px",
                            border: "1px solid #ddd",
                            borderRadius: 6,
                        }}
                    >
                        <div style={{ fontWeight: 600, color: "black" }}>
                            Waiting for an opponent…
                        </div>
                        <div style={{ fontSize: 14 }}>
                            You’re in the lobby as <b>{status.seat}</b>. Share
                            the game link/ID to invite someone.
                        </div>
                    </div>
                )}

                <div>
                    Players: {gameInfo.players[0]} vs{" "}
                    {gameInfo.players[1] || "(waiting)"}
                </div>

                <div>Turn: {gameInfo.turn}</div>

                {!isLive && (
                    <div style={{ margin: "8px 0" }}>
                        Viewing move: <b>{viewMoveNum}</b>{" "}
                        <button onClick={backToLive} style={{ marginLeft: 8 }}>
                            Live
                        </button>
                    </div>
                )}

                <div style={{ marginBottom: 12 }}>
                    <PlayerClock
                        label="⚫ Black"
                        player={2}
                        active={gameInfo.turn === 2}
                    />
                    <PlayerClock
                        label="⚪ White"
                        player={1}
                        active={gameInfo.turn === 1}
                    />
                </div>

                <GoBoardSVG
                    board={boardToShow}
                    interactive={canPlay}
                    onPlay={(r, c) => sendMove(r, c)}
                />

                <button
                    onClick={sendPass}
                    disabled={!canPlay}
                    title={
                        isLobbyWaiting
                            ? "Waiting for opponent to join"
                            : undefined
                    }
                >
                    pass
                </button>
            </div>

            <div style={{ minWidth: 220 }}>
                <GameHistory
                    maxMoveNum={maxMoveNum}
                    selectedMoveNum={viewMoveNum}
                    onSelect={handleSelectMove}
                />
            </div>
        </div>
    );
};

export default Game;
