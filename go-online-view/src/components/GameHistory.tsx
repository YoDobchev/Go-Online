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
        <div style={{ border: "1px solid #ddd", padding: 8, borderRadius: 8 }}>
            <div style={{ fontWeight: 600, marginBottom: 8 }}>History</div>

            <div style={{ display: "flex", flexWrap: "wrap", gap: 6 }}>
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
                            style={{
                                padding: "4px 8px",
                                borderRadius: 6,
                                border: "1px solid #ccc",
                                background: isSelected ? "#9d2c2c" : "#fff",
                                color: "#000000",
                                cursor: "pointer",
                            }}
                        >
                            {label}
                        </button>
                    );
                })}
            </div>
        </div>
    );
}
