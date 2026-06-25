import { SignInForm } from "@/components/features/auth";
import styles from "./page.module.css";

export const SignIn = () => {
    return (
        <div className={`${styles.pageWrapper} ${styles.pageBottomPadding} ${styles.pageCenterContent}`}>
            <SignInForm />
        </div>
    );
};
