import { SignUpForm } from "@/components/features/auth";
import styles from "./pages.module.css";

export const SignUp = () => {
    return (
        <div className={`${styles.pageWrapper} ${styles.pageBottomPadding} ${styles.pageCenterContent}`}>
            <SignUpForm />
        </div>
    );
};
