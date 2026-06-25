import type { ComponentProps, ReactNode } from "react";
import styles from "./AuthFormLayout.module.css";

interface AuthFormLayoutProps {
    title: string;
    onSubmit: ComponentProps<"form">["onSubmit"];
    children: ReactNode;
    note?: ReactNode;
}

export const AuthFormLayout = ({ title, onSubmit, children, note }: AuthFormLayoutProps) => {
    return (
        <form className={styles.form} onSubmit={onSubmit}>
            <h1 className={styles.title}>{title}</h1>
            {children}
            {note && <div className={styles.note}>{note}</div>}
        </form>
    );
};
