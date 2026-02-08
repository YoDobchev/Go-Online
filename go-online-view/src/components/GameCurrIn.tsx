import { useContext } from "react";
import { UserContext } from "../context/UserContext";
import "../styles/GameCurrIn.scss";

const GameCurrIn: React.FC = () => {
    const { user, loading } = useContext(UserContext);

    if (loading || !user) return <div>Checking game status...</div>;

    return (
        <>
            {user.isInGameWithID != "" && (
                <div className="gamecurrin-box">
                    <a href={"/game/" + user.isInGameWithID}>
                        You are in a game! {user.isInGameWithID}
                    </a>
                </div>
            )}
        </>
    );
};

export default GameCurrIn;
