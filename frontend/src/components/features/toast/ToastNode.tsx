import type { Toast } from "@/contexts/ToastContext";
import { Info, Ban, Check, TriangleAlert, X } from "lucide-react";

import styles from "./ToastNode.module.css";
import { useToast } from "@/hooks/useToast";

export const ToastNode = ({ toast } : {toast: Toast}) => {
    const { removeToast } = useToast();

    const renderIcon = () => {
        switch (toast.type) {
            case "success":
                return <Check size={16}/>;
            case "error":
                return <Ban size={16}/>;
            case "warning":
                return <TriangleAlert size={16}/>;
            case "info":
                return <Info size={16}/>;
            default:
                return null;
        }
    };
    
    return (
        <div className={`${styles.toastNode} ${styles[toast.type]}`}>
            <div className={styles.icon}>
                {renderIcon()}
            </div>

            <p className={styles.message}>{toast.message}</p>

            <button
                className={styles.closeButton}
                onClick={() => removeToast(toast.id)}
                aria-label="Close toast"
            >
                <X size={16} />
            </button>

            <div className={styles.progressBar}>
                <div
                    className={styles.progressFill}
                    style={{ animationDuration: `${toast.duration}ms` }}
                />
            </div>
        </div>
    )
}