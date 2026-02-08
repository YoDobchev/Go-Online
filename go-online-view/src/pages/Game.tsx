import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { API_BASE } from "../config";
import GoBoardSVG from "../components/GoBoardSVG";
import type { Board } from "../components/GoBoardSVG";
import GameHistory from "../components/GameHistory";
import PlayerClock from "../components/PlayerClock";
import "../styles/Game.scss";

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

    const getClockView = (player: number) => {
        const clock = clocks[player];
        if (!clock) {
            return {
                remainingMs: null,
                isByoYomi: false,
                byoPeriodsLeft: 0,
                byoPeriodMs: 0,
            };
        }

        let remaining =
            clock.main_remaining_ms > 0
                ? clock.main_remaining_ms
                : clock.byo_remaining_ms;

        if (gameInfo?.turn === player) {
            remaining -= Date.now() - clockReceivedAt;
        }

        return {
            remainingMs: Math.max(0, remaining),
            isByoYomi: clock.main_remaining_ms <= 0,
            byoPeriodsLeft: clock.byo_periods_left,
            byoPeriodMs: clock.byo_remaining_ms,
        };
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
        <div className="game-container">
            <div className="game-left">
                <div className="game-board-wrapper">
                    <GoBoardSVG
                        board={boardToShow}
                        interactive={canPlay}
                        onPlay={(r, c) => sendMove(r, c)}
                    />
                </div>
            </div>

            <div className="game-right">
                <div className="game-status">
                    {status.role === "spectator" && (
                        <div className="game-status__spectator-notice">
                            You are spectating.
                        </div>
                    )}
                </div>

                {isLobbyWaiting && (
                    <div className="game-lobby-waiting">
                        <div className="game-lobby-waiting__title">
                            Waiting for an opponent…
                        </div>
                        <div className="game-lobby-waiting__message">
                            You're in the lobby as <b>{status.seat}</b>. Share
                            the game link/ID to invite someone.
                        </div>
                    </div>
                )}

                <div className="game-players">
                    <strong className="player-black">{gameInfo.players[0]}</strong> vs{" "}
                    <strong className="player-white">{gameInfo.players[1] || "(waiting)"}</strong>
                </div>

              

                <div className="game-clocks">
                    <PlayerClock
                        label="⚫ Black"
                        active={gameInfo.turn === 2}
                        {...getClockView(2)}
                    />
                    <PlayerClock
                        label="⚪ White"
                        active={gameInfo.turn === 1}
                        {...getClockView(1)}
                    />
                </div>
                <div className="game-controls">
                    <div className="game-controls__left">
                        <button
                            onClick={sendPass}
                            disabled={!canPlay}
                            title={
                                isLobbyWaiting
                                    ? "Waiting for opponent to join"
                                    : undefined
                            }
                            className="btn btn--primary"
                        >
                            Pass
                        </button>

                        <button
                            onClick={() => console.log("resign")}
                            disabled={!canPlay}
                            className="btn btn--danger"
                            title={
                                !canPlay
                                    ? "You can only resign on your turn"
                                    : undefined
                            }
                        >
                            Resign
                        </button>
                    </div>

                    <button
                        className="btn btn--icon"
                        onClick={() => console.log("report")}
                        title="Report"
                        aria-label="Report"
                    >
                        <span className="flag" aria-hidden="true">
                            🚩
                        </span>
                    </button>
                </div>

                  {!isLive && (
                    <div className="game-history-notice">
                        <span className="game-history-notice__text">
                            Viewing History!
                        </span>
                        <button onClick={backToLive} className="btlbtn">Back to Live</button>
                    </div>
                )}

                <div className="game-sidebar">
                    <GameHistory
                        maxMoveNum={maxMoveNum}
                        selectedMoveNum={viewMoveNum}
                        onSelect={handleSelectMove}
                    />
                </div>
            </div>
        </div>
    );
};

export default Game;
