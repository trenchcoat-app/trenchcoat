import { useId, type ComponentProps } from "react";
import { useTranslation } from "react-i18next";
import "./Input.css";

interface InputProps extends ComponentProps<"input"> {
    label?: string;
    errors?: (string | undefined)[];
}

export const Input = ({ id, ref, label, errors, ...props }: InputProps) => {
    const { t } = useTranslation(['validation']);

    const generatedId = useId();
    const inputId = id ?? generatedId;

    const activeErrors = errors?.filter(Boolean) ?? [];
    const hasErrors = activeErrors.length > 0;

    return (
        <div className="input-wrapper">
            {label && (
                <label htmlFor={inputId} className="input-label">
                    {label}
                </label>
            )}
            <input id={inputId} className={`input ${hasErrors ? "input-invalid" : ""}`} ref={ref} {...props} />
            <div id={`${inputId}-error`} className="input-error">
                <p>{hasErrors ? t(activeErrors[0] as string) : ""}</p>
            </div>
        </div>
    );
};
