import type { Validator } from "@/types/validator";

export const passwordLengthFieldValidator: Validator = ({ value }: { value: string }) => {
    if (value.length < 8) return "PASSWORD_LENGTH_ERROR";
    return undefined;
};
