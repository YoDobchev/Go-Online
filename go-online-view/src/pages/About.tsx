import Navbar from "../components/Navbar";

const About: React.FC = () => {
    return (
        <div>
            <Navbar/>   
            <h1 style={{margin: "0", padding: "20px"}}>About Go Online</h1>
            <p style={{margin: "0", paddingLeft: "20px"}}>
                Go Online is a platform for playing the ancient board game Go. <br/>
                It was created by Yoan Dobchev and Konstantin Nikolov for the Go programming course at FMI, Sofia university, in 2026.
            </p>
        </div>
    );
};

export default About;