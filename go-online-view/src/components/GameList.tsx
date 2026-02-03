import { useEffect } from "react";
// import { API_BASE } from "../config";

const GameList: React.FC = () => {
    useEffect(() => {
        const fetchGames = async () => {
            // const res = fetch(`${API_BASE}/gameList`, {
            //     credentials: "include",
            // });

        }

        fetchGames();
    });
    return (
        <div>


        </div>
    );
};

export default GameList;
