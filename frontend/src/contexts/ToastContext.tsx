import { createContext, useState, type ReactNode } from "react";

const toastTypes = ["info", "error", "success", "warning"] as const
export type ToastType = (typeof toastTypes)[number];

export type Toast = {
    id: string;
    type: ToastType
    message: string;
}

interface ToastContextType {
    toasts: Toast[];
    addToast: (toast: Omit<Toast, "id">) => void;
    removeToast: (id: string) => void;
}

export const ToastContext = createContext<ToastContextType | null>(null);

export const ToastProvider = ({ children }: { children: ReactNode }) => {
    const [toasts, setToasts] = useState<Toast[]>([]);

    const addToast = (toast: Omit<Toast, "id">) => {
        const id = crypto.randomUUID();
        const newToast: Toast = {
            id,
            ...toast
        }

        setToasts((prev) => [
            ...prev,
            newToast
        ]);

        // setTimeout(() => {
        //     removeToast(id);
        // }, 4000);
    }

    const removeToast = (id: string) => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
    };

    return (
        <ToastContext.Provider value={{ toasts, addToast, removeToast }}>
            {children}
        </ToastContext.Provider>
    )
}