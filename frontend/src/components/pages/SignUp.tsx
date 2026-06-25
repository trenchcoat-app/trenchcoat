import { SignUpForm } from "@/components/features/SignUpForm";
import styles from "@/styles/page-layout.module.css";

export const SignUp = () => {
    return (
        <div className={styles.pageWrapper}>
            <SignUpForm />
        </div>
    );
};
