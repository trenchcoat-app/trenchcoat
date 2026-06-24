import { useId, type ComponentProps } from "react";
import "./Input.css";

interface InputProps extends ComponentProps<"input"> {
    label?: string;
    errors?: string[];
}

export const Input = ({ id, ref, label, errors = [], ...props }: InputProps) => {
    const generatedId = useId();
    const inputId = id ?? generatedId;

    const hasErrors = errors.length > 0;

    return (
        <div className="input-wrapper">
            {label && (
                <label htmlFor={inputId} className="input-label">
                    {label}
                </label>
            )}
            <input id={inputId} className={`input ${hasErrors ? "input-invalid" : ""}`} ref={ref} {...props} />
            <div id={`${inputId}-error`} className="input-error">
                {hasErrors && <p>{errors[0]}</p>}
            </div>
        </div>
    );
};
