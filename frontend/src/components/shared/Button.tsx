import type { ComponentProps, MouseEvent } from "react";
import styles from "./Button.module.css";

interface ButtonProps extends ComponentProps<"button"> {}

export const Button = ({ ref, onClick, ...props }: ButtonProps) => {
    const handleClick = (e: MouseEvent<HTMLButtonElement>) => {
        if (props["aria-disabled"]) {
            e.preventDefault();
            return;
        }
        onClick?.(e);
    };

    return (
        <div className={styles.buttonWrapper}>
            <button className={styles.button} ref={ref} onClick={handleClick} {...props} />
        </div>
    );
};
