import { useContext } from "react";
import { ToastContext } from "@/contexts/ToastContext";

export const useToast = () => {
    const ctx = useContext(ToastContext);

    if (!ctx) {
        throw new Error("useToast must be used within an ToastProvider");
    }
    return ctx;
};
