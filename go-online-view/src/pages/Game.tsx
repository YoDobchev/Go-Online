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
          data: { white_points: number; black_points: number };
      }
    | { type: "error"; data: string };

const Game: React.FC = () => {
    const { gameID } = useParams<{ gameID: string }>();
    const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
    const [status, setStatus] = useState<Status | null>(null);

    const [viewMoveNum, setViewMoveNum] = useState<number | null>(null);
    const [viewBoard, setViewBoard] = useState<Board | null>(null);

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
            } else if (msg.type === "game_ended") {
                console.log(msg.type, msg.data);
                console.log("gameended");
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

    const sendPass = () => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;

        ws.send(
            JSON.stringify({
                type: "play_move",
                data: { row: -1, col: -1 },
            }),
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

    const canPlay = status.role === "player" && isLive && opponentJoined;

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
                        <div style={{ fontWeight: 600, color: "black"}}>
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