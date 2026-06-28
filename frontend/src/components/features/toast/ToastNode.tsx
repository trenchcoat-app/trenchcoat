import { motion } from "framer-motion";
import { Info, Ban, Check, TriangleAlert, X } from "lucide-react";
import type { Toast } from "@/contexts/ToastContext";
import { useToast } from "@/hooks/useToast";

import styles from "./ToastNode.module.css";

export const ToastNode = ({ toast } : {toast: Toast}) => {
    const { removeToast } = useToast();

    const iconSize = 16;
    const renderIcon = () => {
        switch (toast.type) {
            case "success":
                return <Check size={iconSize}/>;
            case "error":
                return <Ban size={iconSize}/>;
            case "warning":
                return <TriangleAlert size={iconSize}/>;
            case "info":
                return <Info size={iconSize}/>;
            default:
                return null;
        }
    };
    
    return (
        <motion.div
            className={`${styles.toastNode} ${styles[toast.type]}`}
            layout
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20, scale: 0.9 }}
            transition={{ duration: 0.35, ease: [0.21, 1.02, 0.73, 1] }}
        >
            <div className={styles.icon}>{renderIcon()}</div>
            
            <p className={styles.message}>{toast.message}</p>

            <button className={styles.closeButton} onClick={() => removeToast(toast.id)}>
                <X size={16} />
            </button>
            
            <div className={styles.progressBar}>
                <div
                    className={styles.progressFill}
                    style={{ animationDuration: `${toast.duration}ms` }}
                />
            </div>
        </motion.div>
    )
}