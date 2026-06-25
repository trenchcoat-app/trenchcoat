import type { ComponentProps, ReactNode } from "react";
import styles from "./AuthFormLayout.module.css";

interface AuthFormLayoutProps {
    onSubmit: ComponentProps<"form">["onSubmit"];
    children: ReactNode;
}

export const AuthFormLayout = ({ onSubmit, children }: AuthFormLayoutProps) => {
    return (
        <form className={styles.form} onSubmit={onSubmit}>
            {children}
        </form>
    );
};
