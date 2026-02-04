import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom";
import { API_BASE } from "../config";
import GoBoardSVG from "../components/GoBoardSVG";
import type { Board } from "../components/GoBoardSVG";


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
