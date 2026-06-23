import type { Validator } from "@/types/validator";

export const requiredFieldValidator: Validator = ({ value }: { value: string; }) => {
    if (value === undefined || value === null || value === "") return "This field is required";
    return undefined;
};
