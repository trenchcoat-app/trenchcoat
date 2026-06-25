import type { ComponentProps } from "react";
import styles from "./Button.module.css";

interface ButtonProps extends ComponentProps<"button"> {}

export const Button = ({ ref, ...props }: ButtonProps) => {
    return (
        <div className={styles.buttonWrapper}>
            <button className={styles.button} ref={ref} {...props} />
        </div>
    );
};
