import { SignInForm } from "@/components/features/SignInForm";
import styles from "./page.module.css";

export const SignIn = () => {
    return (
        <div className={`${styles.pageWrapper} ${styles.pageBottomPadding} ${styles.pageCenterContent}`}>
            <SignInForm />
        </div>
    );
};
