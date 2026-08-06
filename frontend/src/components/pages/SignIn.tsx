import { SignInForm } from "@/components/features/auth";
import styles from "./pages.module.css";

export const SignIn = () => {
    return (
        <div className={`${styles.pageWrapper} ${styles.pageBottomPadding} ${styles.pageCenterContent}`}>
            <SignInForm />
        </div>
    );
};
