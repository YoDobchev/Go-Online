import React from "react";

interface PlayerClockProps {
    label: string;
    active: boolean;
    remainingMs: number | null;
    isByoYomi: boolean;
    byoPeriodsLeft: number;
    byoPeriodMs: number;
}

const formatTime = (ms: number) => {
    const totalSeconds = Math.ceil(ms / 1000);
    const minutes = Math.floor(totalSeconds / 60);
    const seconds = totalSeconds % 60;
    return `${minutes}:${seconds.toString().padStart(2, "0")}`;
};

const PlayerClock: React.FC<PlayerClockProps> = ({
    label,
    active,
    remainingMs,
    isByoYomi,
    byoPeriodsLeft,
    byoPeriodMs,
}) => {
    return (
        <div
            style={{
                padding: "6px 10px",
                borderRadius: 6,
                border: "1px solid #ddd",
                background: active ? "#eef6ff" : "#fafafa",
                marginBottom: 6,
            }}
        >
            <div style={{ fontWeight: 600 }}>{label}</div>

            <div style={{ fontSize: 18 }}>
                {remainingMs !== null ? formatTime(remainingMs) : "--:--"}
            </div>

            {isByoYomi && (
                <div style={{ fontSize: 12, color: "#555" }}>
                    Byo-yomi: {byoPeriodsLeft} × {formatTime(byoPeriodMs)}
                </div>
            )}
        </div>
    );
};

export default PlayerClock;
