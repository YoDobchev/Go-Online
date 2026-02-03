import { useEffect, useState } from "react";
import { API_BASE } from "../config";

interface GameStatus {
    id: string;
}

const GameCurrIn: React.FC = () => {
    const [gameId, setGameId] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    useEffect(() => {
        const fetchGameStatus = async () => {
            try {
                const res = await fetch(`${API_BASE}/game/`, {
                    credentials: "include",
                });

                const data: GameStatus = await res.json();
                setGameId(data.id == "" ? null : data.id);
            } catch (err) {
                console.error("Failed to fetch game status", err);
                setGameId(null);
            } finally {
                setLoading(false);
            }
        };
        fetchGameStatus();
    }, []);

    const createNewGame = async () => {
        try {
            const res = await fetch(`${API_BASE}/game/`, {
                method: "POST",
                credentials: "include",
            });
            if (!res.ok) {
                throw new Error("Failed to create new game");
            }
            const data: GameStatus = await res.json();
            setGameId(data.id);
        } catch (err) {
            console.error("Failed to create new game", err);
        }
    };


    if (loading) return <div>Checking game status...</div>;

    return (
        <div>
            {gameId ? (
                <div>You are in a game! ID: {gameId}</div>
            ) : (
                <button onClick={createNewGame}>Create new game</button>
            )}
        </div>
    );
};

export default GameCurrIn;
