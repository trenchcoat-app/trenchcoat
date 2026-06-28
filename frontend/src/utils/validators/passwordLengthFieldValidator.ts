import type { Validator } from "@/utils/validators/validator";

export const passwordLengthFieldValidator: Validator = ({ value }: { value: string }) => {
    if (value.length < 8) return "PASSWORD_LENGTH_ERROR";
    return undefined;
};
