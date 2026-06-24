import type { Validator } from "@/types/validator";

export const passwordLengthFieldValidator: Validator = ({ value }: { value: string }) => {
    if (value.length < 8) return "Password must be at least 8 characters";
    return undefined;
};
