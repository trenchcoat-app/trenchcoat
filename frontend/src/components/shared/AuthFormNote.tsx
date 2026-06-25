import type { ReactNode } from "react";
import styles from "./AuthFormNote.module.css";

interface AuthFormNoteProps {
    children: ReactNode;
}

export const AuthFormNote = ({ children }: AuthFormNoteProps) => {
    return <div className={styles.note}><p>{children}</p></div>;
};
