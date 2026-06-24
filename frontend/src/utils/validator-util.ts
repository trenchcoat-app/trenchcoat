import type { Validator, ValidatorContext } from "@/types/validator";

/**
 * Composes multiple validators into a single validator function.
 *
 * Each validator is executed with the same validation context, and any returned
 * error messages (truthy values) are collected into an array.
 *
 * @param validators - A list of validator functions to run in sequence.
 * @returns A single validator function that returns an array of error messages.
 *          If no validators return errors, the array will be empty.
 */
export const composeValidators =
    (...validators: Validator[]) =>
    (context: ValidatorContext) => {
        const errors = validators.map((validator) => validator(context)).filter(Boolean);

        return errors;
    };
