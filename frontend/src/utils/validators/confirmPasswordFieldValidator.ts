import type { Validator } from "@/utils/validators/validator";

export const confirmPasswordFieldValidator: Validator<string> = ({ value, fieldApi }) => {
    const password = fieldApi.form.getFieldValue("password");
    if (value !== password) return "PASSWORDS_MUST_MATCH_ERROR";
    return undefined;
};