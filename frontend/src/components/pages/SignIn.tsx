import { SignInForm } from "@/components/features/SignInForm";
import styles from "@/styles/page-layout.module.css";

export const SignIn = () => {
    return (
        <div className={styles.pageWrapper}>
            <SignInForm />
        </div>
    );
};
