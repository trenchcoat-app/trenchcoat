import type { AnyFieldApi } from "@tanstack/react-form";

export type ValidatorContext = {
    value: string;
    fieldApi: AnyFieldApi;
};

export type Validator = (context: ValidatorContext) => string | undefined;
