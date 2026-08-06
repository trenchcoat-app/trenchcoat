import type { AnyFieldApi } from "@tanstack/react-form";

export type ValidatorContext<TValue = unknown> = {
    value: TValue;
    fieldApi: AnyFieldApi;
};

export type Validator<TValue = unknown> = (context: ValidatorContext<TValue>) => string | undefined;