import { useContext } from "react";
import { useNavigate } from "react-router-dom";
import { API_BASE } from "../config";
import { UserContext } from "../context/UserContext";
import "../styles/GameCurrIn.scss";

interface GameStatus {
    id: string;
}

const GameCurrIn: React.FC = () => {
    const { user, loading } = useContext(UserContext);

    const navigate = useNavigate();

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
            navigate(`/game/${data.id}`, { replace: true });
        } catch (err) {
            console.error("Failed to create new game", err);
        }
    };

    if (loading || !user) return <div>Checking game status...</div>;

    return (
        <>
            {user.isInGameWithID != "" ? (
                <div className="gamecurrin-box">
                    <a href={"/game/" + user.isInGameWithID}>You are in a game! {user.isInGameWithID}</a>
                </div>
            ) : (
                <button onClick={createNewGame}>Create new game</button>
            )}
        </>
    );
};

export default GameCurrIn;
