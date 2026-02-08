export type Stone = 0 | 1 | 2;

export interface Board {
    size: number;
    squares: number[][];
}

export function GoBoardSVG({
    board,
    onPlay,
    interactive = false,
}: {
    board?: Board;
    onPlay?: (row: number, col: number) => void;
    interactive?: boolean;
}) {
    if (!board) {
        return <div>Loading board...</div>;
    }
    const N = board.size;
    const padding = 24;
    const cell = 28;
    const w = padding * 2 + cell * (N - 1);
    const h = padding * 2 + cell * (N - 1);

    const x = (col: number) => padding + col * cell;
    const y = (row: number) => padding + row * cell;

    const starCoords = N === 19 ? [3, 9, 15] : [];
    const starPoints =
        starCoords.length > 0
            ? starCoords.flatMap((r) => starCoords.map((c) => ({ r, c })))
            : [];

    const handleClick = (row: number, col: number) => {
        if (!interactive) return;
        if (!onPlay) return;
        onPlay(row, col);
    };

    return (
        <svg
            width={w}
            height={h}
            viewBox={`0 0 ${w} ${h}`}
            style={{
                display: "block",
                maxWidth: "100%",
                height: "auto",
                userSelect: "none",
            }}
        >
            <rect x={0} y={0} width={w} height={h} fill="#deb887" />

            {Array.from({ length: N }).map((_, i) => (
                <g key={i}>
                    <line
                        x1={x(i)}
                        y1={y(0)}
                        x2={x(i)}
                        y2={y(N - 1)}
                        stroke="rgba(0,0,0,0.75)"
                        strokeWidth={1}
                    />
                    <line
                        x1={x(0)}
                        y1={y(i)}
                        x2={x(N - 1)}
                        y2={y(i)}
                        stroke="rgba(0,0,0,0.75)"
                        strokeWidth={1}
                    />
                </g>
            ))}

            {starPoints.map((p, idx) => (
                <circle
                    key={idx}
                    cx={x(p.c)}
                    cy={y(p.r)}
                    r={3}
                    fill="rgba(0,0,0,0.85)"
                />
            ))}

            {Array.from({ length: N }).map((_, r) =>
                Array.from({ length: N }).map((__, c) => (
                    <circle
                        key={`hit-${r}-${c}`}
                        cx={x(c)}
                        cy={y(r)}
                        r={cell * 0.45}
                        fill="transparent"
                        style={{ cursor: interactive ? "pointer" : "default" }}
                        onClick={() => handleClick(r, c)}
                    />
                )),
            )}

            {board.squares.map((row, r) =>
                row.map((v, c) => {
                    const stone = v as Stone;
                    if (stone === 0) return null;

                    const isBlack = stone === 2;
                    return (
                        <g key={`stone-${r}-${c}`}>
                            <circle
                                cx={x(c) + 1.2}
                                cy={y(r) + 1.2}
                                r={cell * 0.42}
                                fill="rgba(0,0,0,0.25)"
                            />
                            <circle
                                cx={x(c)}
                                cy={y(r)}
                                r={cell * 0.42}
                                fill={isBlack ? "#111" : "#f5f5f5"}
                                stroke={isBlack ? "#000" : "#bbb"}
                                strokeWidth={1}
                            />
                            <circle
                                cx={x(c) - cell * 0.12}
                                cy={y(r) - cell * 0.12}
                                r={cell * 0.12}
                                fill={
                                    isBlack
                                        ? "rgba(255,255,255,0.12)"
                                        : "rgba(255,255,255,0.55)"
                                }
                            />
                        </g>
                    );
                }),
            )}
        </svg>
    );
}

export default GoBoardSVG;
