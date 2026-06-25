import styles from "./AuthFormTitle.module.css";

interface AuthFormTitleProps {
    title: string;
}

export const AuthFormTitle = ({ title }: AuthFormTitleProps) => {
    return <h1 className={styles.title}>{title}</h1>;
};
