import { useParams } from "react-router-dom";
import { API_BASE } from "../config";
import { useEffect, useState } from "react";

interface Board {
    size: number;
    squares: number[][];
}
interface GameInfo {
    players: [string, string];
    turn: number;
    board: Board;
}

const Game: React.FC = () => {
    const { gameID } = useParams<{ gameID: string }>();

    const [gameInfo, setGameInfo] = useState<GameInfo | null>(null);

    useEffect(() => {
        if (!gameID) return;
        const fetchGameData = async () => {
            const res = await fetch(`${API_BASE}/game/${gameID}`, {
                credentials: "include",
            });
            if (res.ok) {
                const data: GameInfo = await res.json();
                setGameInfo(data);
            }
        };

        fetchGameData();
    }, [gameID]);

    return (
        <div>
            gameInfo: {gameInfo ? JSON.stringify(gameInfo) : "Loading..."}
        </div>
    );
};

export default Game;
