import type { AnyFieldApi } from "@tanstack/react-form";
import type { Validator } from "@/types/validator";

export const confirmPasswordFieldValidator: Validator = ({ value, fieldApi }: { value: string; fieldApi: AnyFieldApi }) => {
    const password = fieldApi.form.getFieldValue("password");
    if (value !== password) return "PASSWORDS_MUST_MATCH_ERROR";
    return undefined;
};
