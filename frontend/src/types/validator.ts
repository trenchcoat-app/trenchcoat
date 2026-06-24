import type { AnyFieldApi } from "@tanstack/react-form";

export type ValidatorContext = {
    value: any;
    fieldApi: AnyFieldApi;
};

export type Validator = (context: ValidatorContext) => string | undefined;
