import Navbar from "../components/Navbar";
import GameCurrIn from "../components/GameCurrIn";
import GameList from "../components/GameList";
import "../styles/Home.scss";

const Home: React.FC = () => {
    return (
        <div className="homePage">
            <Navbar />
            <GameCurrIn />

            <GameList />
        </div>
    );
};

export default Home;
