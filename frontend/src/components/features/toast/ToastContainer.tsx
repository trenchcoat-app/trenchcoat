import { useToast } from "@/hooks/useToast"
import { ToastNode } from "@/components/features/toast/ToastNode";

import styles from "./ToastContainer.module.css";

export const ToastContainer = () => {
    const { toasts } = useToast();
    return (
        <div className={styles.toastContainer}>
            {toasts.map((toast) => (
                <ToastNode key={toast.id} toast={toast}/>
            ))}
        </div>
    )
}