import { useRef, useState, type ReactNode } from "react";
import { ToastContext, type Toast } from "./ToastContext";

const TOAST_DURATION = 4000;

type ToastTimer = {
    timeoutId: ReturnType<typeof setTimeout>;
    startedAt: number;
    remaining: number;
};

export const ToastProvider = ({ children }: { children: ReactNode }) => {
    const [toasts, setToasts] = useState<Toast[]>([]);
    const timers = useRef<Map<string, ToastTimer>>(new Map());

    const removeToast = (id: string) => {
        const timer = timers.current.get(id);
        if (timer) {
            clearTimeout(timer.timeoutId);
            timers.current.delete(id);
        }

        setToasts((prev) => prev.filter((t) => t.id !== id));
    };

    const addToast = (toast: Omit<Toast, "id" | "duration">) => {
        const id = crypto.randomUUID();
        const newToast: Toast = {
            id,
            ...toast,
            duration: TOAST_DURATION,
        };

        setToasts((prev) => [...prev, newToast]);

        const timeoutId = setTimeout(() => removeToast(id), TOAST_DURATION);

        timers.current.set(id, {
            timeoutId,
            startedAt: Date.now(),
            remaining: TOAST_DURATION,
        });
    };

    const pauseToast = (id: string) => {
        const timer = timers.current.get(id);
        if (!timer) return;

        clearTimeout(timer.timeoutId);

        timers.current.set(id, {
            ...timer,
            remaining: timer.remaining - (Date.now() - timer.startedAt),
        });
    };

    const resumeToast = (id: string) => {
        const timer = timers.current.get(id);
        if (!timer) return;

        const timeoutId = setTimeout(() => removeToast(id), timer.remaining);

        timers.current.set(id, {
            timeoutId,
            startedAt: Date.now(),
            remaining: timer.remaining,
        });
    };

    return (
        <ToastContext.Provider
            value={{
                toasts,
                addToast,
                removeToast,
                pauseToast,
                resumeToast,
            }}
        >
            {children}
        </ToastContext.Provider>
    );
};