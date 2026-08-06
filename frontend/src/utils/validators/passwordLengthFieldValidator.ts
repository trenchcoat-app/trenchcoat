import type { Validator } from "@/utils/validators/validator";

export const passwordLengthFieldValidator: Validator<string> = ({ value }) => {
    if (!value || value.length < 8) return "PASSWORD_LENGTH_ERROR";
    return undefined;
};