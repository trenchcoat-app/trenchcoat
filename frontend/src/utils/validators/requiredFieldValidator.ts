import type { Validator } from "@/utils/validators/validator";

export const requiredFieldValidator: Validator = ({ value }: { value: any }) => {
    if (value === undefined || value === null || value === "") return "REQUIRED_FIELD_ERROR";
    return undefined;
};
