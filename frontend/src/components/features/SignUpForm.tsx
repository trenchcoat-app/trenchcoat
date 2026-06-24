import { useMutation } from "@tanstack/react-query";
import { useForm } from "@tanstack/react-form";
import { signUpMutation } from "@/api/@tanstack/react-query.gen";
import { Input, Button } from "@/components/shared";
import type { SignUpBody } from "@/api/types.gen";
import { requiredFieldValidator, confirmPasswordFieldValidator, passwordLengthFieldValidator } from "@/utils/validators";
import { composeValidators } from "@/utils/validator-util";
import { useTranslation } from "react-i18next";

export const SignUpForm = () => {
    const { t } = useTranslation();
    const mutation = useMutation(signUpMutation());

    const defaultValues: SignUpBody & { confirmPassword: string } = {
        email: "",
        password: "",
        displayName: "",
        confirmPassword: "",
    };

    const form = useForm({
        defaultValues,
        onSubmit: ({ value }) => {
            const { confirmPassword, ...body } = value;
            mutation.mutate({ body });
        },
    });

    return (
        <form
            onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
            }}
        >
            <form.Field
                name="email"
                validators={{
                    onChange: requiredFieldValidator,
                }}
            >
                {(field) => (
                    <Input
                        name={field.name}
                        autoComplete={field.name}
                        type="email"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label={t("auth:EMAIL")}
                        placeholder={t("auth:EMAIL_PLACEHOLDER")}
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Field
                name="displayName"
                validators={{
                    onChange: requiredFieldValidator,
                }}
            >
                {(field) => (
                    <Input
                        name={field.name}
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label={t("auth:DISPLAY_NAME")}
                        placeholder={t("auth:DISPLAY_NAME")}
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Field
                name="password"
                validators={{
                    onChange: composeValidators(requiredFieldValidator, passwordLengthFieldValidator),
                }}
            >
                {(field) => (
                    <Input
                        name={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label={t("auth:PASSWORD")}
                        placeholder={t("auth:PASSWORD")}
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Field
                name="confirmPassword"
                validators={{
                    onChange: confirmPasswordFieldValidator,
                }}
            >
                {(field) => (
                    <Input
                        name={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label={t("auth:CONFIRM_PASSWORD")}
                        placeholder={t("auth:CONFIRM_PASSWORD")}
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
                children={([canSubmit, isSubmitting]) => (
                    <Button type="submit" disabled={!canSubmit}>
                        {isSubmitting ? "..." : t("auth:SIGNUP")}
                    </Button>
                )}
            />
        </form>
    );
};
