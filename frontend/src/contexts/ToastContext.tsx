import { createContext, useState, type ReactNode } from "react";

const TOAST_DURATION = 4000;

const toastTypes = ["info", "error", "success", "warning"] as const;
export type ToastType = (typeof toastTypes)[number];

export type Toast = {
    id: string;
    type: ToastType;
    message: string;
    duration: number;
};

interface ToastContextType {
    toasts: Toast[];
    addToast: (toast: Omit<Toast, "id" | "duration">) => void;
    removeToast: (id: string) => void;
}

export const ToastContext = createContext<ToastContextType | null>(null);

export const ToastProvider = ({ children }: { children: ReactNode }) => {
    const [toasts, setToasts] = useState<Toast[]>([]);

    const addToast = (toast: Omit<Toast, "id" | "duration">) => {
        const id = crypto.randomUUID();
        const newToast: Toast = {
            id,
            ...toast,
            duration: TOAST_DURATION,
        };

        setToasts((prev) => [...prev, newToast]);

        setTimeout(() => {
            removeToast(id);
        }, newToast.duration);
    };

    const removeToast = (id: string) => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
    };

    return <ToastContext.Provider value={{ toasts, addToast, removeToast }}>{children}</ToastContext.Provider>;
};
