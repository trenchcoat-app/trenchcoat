import type { Validator, ValidatorContext } from "@/types/validator";

export const composeValidators = (...validators: Validator[]) =>
    (context: ValidatorContext) => {
        const errors = validators
            .map((validator) => validator(context))
            .filter(Boolean);

        return errors;
    };