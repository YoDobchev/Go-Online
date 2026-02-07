import Navbar from "../components/Navbar";
import GameCurrIn from "../components/GameCurrIn";
import GameList from "../components/GameList";

const Home: React.FC = () => {
    return (
        <div>
            <Navbar />
            <GameCurrIn />
            <br />
            <br />
            <br />
            <br />
            <br />
            <br />
            <GameList />
        </div>
    );
};

export default Home;
