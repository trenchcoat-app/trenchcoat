import { useId, type ComponentProps } from "react";
import styles from "./Input.module.css";

interface InputProps extends ComponentProps<"input"> {
    label?: string;
    errors?: string[];
}

export const Input = ({ id, ref, label, errors = [], ...props }: InputProps) => {
    const generatedId = useId();
    const inputId = id ?? generatedId;

    const hasErrors = errors.length > 0;

    return (
        <div className={styles.inputWrapper}>
            {label && (
                <label htmlFor={inputId} className={styles.inputLabel}>
                    {label}
                </label>
            )}
            <input
                id={inputId}
                className={`${styles.input} ${hasErrors ? styles.inputInvalid : ""}`}
                ref={ref}
                aria-describedby={hasErrors ? `${inputId}-error` : undefined}
                aria-invalid={hasErrors}
                {...props}
            />
            {hasErrors && (
                <div id={`${inputId}-error`} className={styles.inputError} role="alert">
                    <p>{errors[0]}</p>
                </div>
            )}
        </div>
    );
};
