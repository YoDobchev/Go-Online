import { useEffect, useRef, useState, useCallback, useContext } from "react";
import { useParams } from "react-router-dom";
import { API_BASE } from "../config";
import GoBoardSVG from "../components/GoBoardSVG";
import type { Board } from "../components/GoBoardSVG";
import GameHistory from "../components/GameHistory";
import PlayerClock from "../components/PlayerClock";
import "../styles/Game.scss";
import Navbar from "../components/Navbar";
import BottomCenMsg from "../components/BottomCenMsg";
import { UserContext } from "../context/UserContext";

interface GameInfo {
    players: [string, string];
    turn: number;
    board: Board;
    moveNum: number;
}

interface MoveResp {
    board: Board;
    blackWinProb: number;
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
          data: {
              players: [string, string];
              white_points: number;
              black_points: number;
              winner: number;
              reason: string;
              moveNum: number;
          };
      }
    | { type: "error"; data: string }
    | { type: "clock_update"; data: ClockSnapshot };

const Game: React.FC = () => {
    const { user } = useContext(UserContext);

    const { gameID } = useParams<{ gameID: string }>();
    const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);
    const [status, setStatus] = useState<Status | null>(null);

    const [viewMoveNum, setViewMoveNum] = useState<number | null>(null);
    const [viewBoard, setViewBoard] = useState<Board | null>(null);
    const [blackWinProb, setBlackWinProb] = useState<number | null>(null);

    const [clocks, setClocks] = useState<Record<number, ClockSnapshot>>({});
    const [clockReceivedAt, setClockReceivedAt] = useState<number>(Date.now());

    const wsRef = useRef<WebSocket | null>(null);

    const [show, setShow] = useState(false);
    const [msg, setMsg] = useState("");
    const handleCloseMsg = useCallback(() => setShow(false), []);

    const [gameEndedInfo, setGameEndedInfo] = useState<{
        players: [string, string];
        white_points: number;
        black_points: number;
        winner: number;
        reason: string;
        moveNum: number;
    } | null>(null);

    const handleSelectMove = useCallback(
        async (moveNum: number) => {
            console.log("Selecting move", moveNum);
            if (!gameID) throw new Error("No gameID");
            setViewMoveNum(moveNum);
            try {
                const res = await fetch(
                    `${API_BASE}/game/${gameID}/state?moveNum=${moveNum}`,
                );
                if (!res.ok) throw new Error(await res.text());
                const b = (await res.json()) as MoveResp;
                setViewBoard(b.board);
                setBlackWinProb(b.blackWinProb);
            } catch (e) {
                console.error("Failed to load snapshot:", e);
            }
        },
        [gameID],
    );

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
                setGameEndedInfo(msg.data);
                console.log("Game ended:", msg.data);
                handleSelectMove(msg.data.moveNum);
            } else if (msg.type === "clock_update") {
                setClocks((prev) => ({
                    ...prev,
                    [msg.data.player]: msg.data,
                }));
                setClockReceivedAt(Date.now());
            } else if (msg.type === "error") {
                console.error("Server error:", msg.data);
                setMsg(msg.data);
                setShow(true);
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
    }, [gameID, handleSelectMove]);

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

    const sendResign = () => {
        const ws = wsRef.current;
        if (!ws || ws.readyState !== WebSocket.OPEN) return;

        ws.send(
            JSON.stringify({
                type: "play_resign",
            }),
        );
    };

    const sendReport = async () => {
        const res = await fetch(`${API_BASE}/reports`, {
            method: "POST",
            credentials: "include",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username: user?.username, game_id: gameID }),
        });
        if (!res.ok) {
            setMsg("Failed to report game");
            setShow(true);
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

    const defaultStatus = status ?? {
        role: "spectator",
        seat: "spectator",
    };

    if (!gameEndedInfo && !gameInfo) return <div>Loading game info...</div>;

    const maxMoveNum = gameEndedInfo?.moveNum ?? gameInfo?.moveNum ?? 0;
    const boardToShow = viewBoard ?? gameInfo?.board;
    const isLive = viewMoveNum === null;

    const opponentJoined = Boolean(gameInfo?.players[1]);
    const isLobbyWaiting = defaultStatus.role === "player" && !opponentJoined;

    const canPlay =
        !gameEndedInfo &&
        defaultStatus.role === "player" &&
        isLive &&
        opponentJoined;

    return (
        <>
            <Navbar></Navbar>

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
                        {defaultStatus.role === "spectator" &&
                            !gameEndedInfo && (
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
                                You're in the lobby as{" "}
                                <b>{defaultStatus.seat}</b>. Share the game
                                link/ID to invite someone.
                            </div>
                        </div>
                    )}

                    {gameInfo && !gameEndedInfo && (
                        <>
                            <div className="game-players">
                                <strong className="player-black">
                                    {gameInfo.players[0]}
                                </strong>{" "}
                                vs{" "}
                                <strong className="player-white">
                                    {gameInfo.players[1] || "(waiting)"}
                                </strong>
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
                        </>
                    )}

                    {gameEndedInfo && (
                        <div className="game-status__ended-notice">
                            <div>
                                <strong className="player-black">
                                    {gameEndedInfo.players[0]}
                                </strong>{" "}
                                vs{" "}
                                <strong className="player-white">
                                    {gameEndedInfo.players[1]}
                                </strong>
                            </div>

                            {gameEndedInfo.reason === "resignation" ||
                            gameEndedInfo.reason === "timeout" ? (
                                <span>
                                    {gameEndedInfo.winner === 1
                                        ? "White"
                                        : "Black"}{" "}
                                    won by{" "}
                                    {gameEndedInfo.reason === "resignation"
                                        ? "resignation"
                                        : "timeout"}
                                    .
                                </span>
                            ) : (
                                <span>
                                    Game ended. Black:{" "}
                                    {gameEndedInfo.black_points} | White:{" "}
                                    {gameEndedInfo.white_points}. Winner:{" "}
                                    {gameEndedInfo.winner === 1
                                        ? "Black"
                                        : "White"}
                                </span>
                            )}
                        </div>
                    )}

                    {defaultStatus.role === "player" && !gameEndedInfo && (
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
                                    onClick={() => sendResign()}
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
                                onClick={() => sendReport()}
                                title="Report"
                                aria-label="Report"
                            >
                                <span className="flag" aria-hidden="true">
                                    🚩
                                </span>
                            </button>
                        </div>
                    )}

                    {!isLive && !gameEndedInfo && (
                        <div className="game-history-notice">
                            <span className="game-history-notice__text">
                                Viewing History!
                            </span>
                            <button onClick={backToLive} className="btlbtn">
                                Back to Live
                            </button>
                        </div>
                    )}

                    {blackWinProb != null && blackWinProb != -1 && (
                        <div className="winprob">
                            <div className="winprob__header">
                                <div className="winprob__title">Prediction</div>
                                <div className="winprob__nums">
                                    <span className="winprob__num winprob__num--black">
                                        ⚫ {Math.round(blackWinProb * 100)}%
                                    </span>
                                    <span className="winprob__num winprob__num--white">
                                        ⚪{" "}
                                        {Math.round((1 - blackWinProb) * 100)}%
                                    </span>
                                </div>
                            </div>

                            <div
                                className="winprob__bar"
                                role="img"
                                aria-label={`Predicted win chance: Black ${Math.round(
                                    blackWinProb * 100,
                                )}%, White ${Math.round((1 - blackWinProb) * 100)}%`}
                            >
                                <div
                                    className="winprob__seg winprob__seg--black"
                                    style={{
                                        width: `${Math.max(0, Math.min(1, blackWinProb)) * 100}%`,
                                    }}
                                />
                                <div
                                    className="winprob__mid"
                                    aria-hidden="true"
                                />
                                <div
                                    className="winprob__seg winprob__seg--white"
                                    style={{
                                        width: `${Math.max(0, Math.min(1, 1 - blackWinProb)) * 100}%`,
                                    }}
                                />
                            </div>
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

            <BottomCenMsg
                visible={show}
                message={msg}
                backgroundColor="#e90e0e"
                textColor="#f7f7f7"
                timeAfterFadeMs={1000}
                onClose={handleCloseMsg}
            />
        </>
    );
};

export default Game;
