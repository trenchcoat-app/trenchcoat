import { createContext, useState, useRef, type ReactNode } from "react";

const TOAST_DURATION = 4000;

const toastTypes = ["info", "error", "success", "warning"] as const;
export type ToastType = (typeof toastTypes)[number];

export type Toast = {
    id: string;
    type: ToastType;
    message: string;
    duration: number;
};

type ToastTimer = {
    timeoutId: ReturnType<typeof setTimeout>;
    startedAt: number;
    remaining: number;
};

interface ToastContextType {
    toasts: Toast[];
    addToast: (toast: Omit<Toast, "id" | "duration">) => void;
    removeToast: (id: string) => void;
    pauseToast: (id: string) => void;
    resumeToast: (id: string) => void;
}

export const ToastContext = createContext<ToastContextType | null>(null);

export const ToastProvider = ({ children }: { children: ReactNode }) => {
    const [toasts, setToasts] = useState<Toast[]>([]);
    const timers = useRef<Map<string, ToastTimer>>(new Map());

    const removeToast = (id: string) => {
        timers.current.delete(id);

        setToasts((prev) => prev.filter((t) => t.id !== id));
    };

    const addToast = (toast: Omit<Toast, "id" | "duration">) => {
        const id = crypto.randomUUID();
        const newToast: Toast = { id, ...toast, duration: TOAST_DURATION };
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

    return <ToastContext.Provider value={{ toasts, addToast, removeToast, pauseToast, resumeToast }}>{children}</ToastContext.Provider>;
};
