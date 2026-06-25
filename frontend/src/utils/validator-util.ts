import type { Validator, ValidatorContext } from "@/types/validator";
import type { TFunction } from "i18next";

// Runs multiple validators and collects their results.
export const composeValidators =
    (...validators: Validator[]) =>
    (context: ValidatorContext): string[] | undefined => {
        const results = validators.map((validator) => validator(context));
        const errors = extractErrors(results);

        return errors.length ? errors : undefined;
    };

// Narrows a mixed array of validation results down to only the error strings.
export const extractErrors = (errors: (string | undefined)[]): string[] => errors.filter((e): e is string => typeof e === "string");

// Translates validation error keys into human-readable messages using the "validation" i18n namespace via i18next's TFunction.
export const localizeErrors = (errors: string[], t: TFunction): string[] => {
    return errors.map((key) => t(`validation:${key}`));
};

// Wrapper function to call extractErrors and localizeErrors in combination
export const extractAndLocalizeErrors = (errors: (string | undefined)[], t: TFunction): string[] => {
    return localizeErrors(extractErrors(errors), t);
};
