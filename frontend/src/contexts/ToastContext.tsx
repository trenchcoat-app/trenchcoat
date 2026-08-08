import { createContext } from "react";

export type ToastType = "info" | "error" | "success" | "warning";

export type Toast = {
    id: string;
    type: ToastType;
    message: string;
    duration: number;
};

export interface ToastContextType {
    toasts: Toast[];
    addToast: (toast: Omit<Toast, "id" | "duration">) => void;
    removeToast: (id: string) => void;
    pauseToast: (id: string) => void;
    resumeToast: (id: string) => void;
}

export const ToastContext = createContext<ToastContextType | null>(null);