import { SignUpForm } from "@/components/features/SignUpForm";
import styles from "./page.module.css";

export const SignUp = () => {
    return (
        <div className={`${styles.pageWrapper} ${styles.pageBottomPadding} ${styles.pageCenterContent}`}>
            <SignUpForm />
        </div>
    );
};
