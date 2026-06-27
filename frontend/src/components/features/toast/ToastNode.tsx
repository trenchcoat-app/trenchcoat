import type { Toast } from "@/contexts/ToastContext";

import styles from "./ToastNode.module.css";

export const ToastNode = ({ toast } : {toast: Toast}) => {
    return (
        <div className={`${styles.toastNode} ${styles[toast.type]}`}>
            <p className={styles.message}>{toast.message}</p>
        </div>
    )
}