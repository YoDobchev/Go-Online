import { useEffect, useMemo, useState } from "react";

export type BottomCenMsgProps = {
    message: string;
    backgroundColor?: string;
    textColor?: string;
    timeAfterFadeMs?: number;
    visible?: boolean;
    onClose?: () => void;
};

export default function BottomCenMsg({
    message,
    backgroundColor = "#111827",
    textColor = "#ffffff",
    timeAfterFadeMs = 2500,
    visible = true,
    onClose,
}: BottomCenMsgProps) {
    const FADE_MS = 250;

    const [mounted, setMounted] = useState(false);
    const [fading, setFading] = useState(false);

    const hasMessage = useMemo(
        () => typeof message === "string" && message.trim().length > 0,
        [message],
    );

    useEffect(() => {
        if (!visible || !hasMessage) {
            setMounted(false);
            setFading(false);
            return;
        }

        setMounted(true);
        setFading(false);

        const t1 = window.setTimeout(
            () => setFading(true),
            Math.max(0, timeAfterFadeMs),
        );
        const t2 = window.setTimeout(
            () => {
                setMounted(false);
                setFading(false);
                onClose?.();
            },
            Math.max(0, timeAfterFadeMs) + FADE_MS,
        );

        return () => {
            window.clearTimeout(t1);
            window.clearTimeout(t2);
        };
    }, [visible, hasMessage, timeAfterFadeMs, onClose]);

    if (!mounted) return null;

    return (
        <div
            style={{
                position: "fixed",
                left: "50%",
                bottom: "24px",
                transform: "translateX(-50%)",
                zIndex: 9999,
                pointerEvents: "none",
            }}
            aria-live="polite"
            aria-atomic="true"
        >
            <div
                style={{
                    background: backgroundColor,
                    color: textColor,
                    padding: "10px 14px",
                    borderRadius: "12px",
                    boxShadow: "0 10px 25px rgba(0,0,0,0.25)",
                    maxWidth: "min(92vw, 520px)",
                    fontSize: "14px",
                    lineHeight: 1.35,
                    opacity: fading ? 0 : 1,
                    transition: `opacity ${FADE_MS}ms ease`,
                }}
            >
                {message}
            </div>
        </div>
    );
}
