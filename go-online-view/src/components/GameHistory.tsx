import "../styles/GameHistory.scss";
export default function GameHistory({
    maxMoveNum,
    selectedMoveNum,
    onSelect,
}: {
    maxMoveNum: number;
    selectedMoveNum: number | null;
    onSelect: (moveNum: number) => void;
}) {
    const moves = Array.from({ length: maxMoveNum + 1 }, (_, i) => i);

    return (
        <div className="game-history">
            <div className="game-history__title">History</div>

            <div className="game-history__grid">
                {moves.map((mn) => {
                    const isSelected = selectedMoveNum === mn;
                    const label =
                        mn === 0
                            ? "Start"
                            : `${mn} ${mn % 2 === 1 ? "B" : "W"}`;

                    return (
                        <button
                            key={mn}
                            onClick={() => onSelect(mn)}
                            className={`game-history__btn ${isSelected ? "is-selected" : ""}`}
                            type="button"
                        >
                            {label}
                        </button>
                    );
                })}
            </div>
        </div>
    );
}
