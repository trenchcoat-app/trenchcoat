import { useMutation } from "@tanstack/react-query";
import { useForm } from "@tanstack/react-form";
import { Link } from "@tanstack/react-router";
import { useTranslation } from "react-i18next";
import { signUpMutation } from "@/api/@tanstack/react-query.gen";
import type { SignUpBody } from "@/api/types.gen";
import { Input, Button } from "@/components/shared";
import { requiredFieldValidator, confirmPasswordFieldValidator, passwordLengthFieldValidator } from "@/utils/validators";
import { composeValidators, extractAndLocalizeErrors } from "@/utils/validator-util";
import "./SignUpForm.css";

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
            className="form"
        >   
            <h1 className="form-title">{t("auth:SIGNUP_TITLE")}</h1>

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
                        errors={extractAndLocalizeErrors(field.state.meta.errors, t)}
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
                        errors={extractAndLocalizeErrors(field.state.meta.errors, t)}
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
                        errors={extractAndLocalizeErrors(field.state.meta.errors, t)}
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
                        errors={extractAndLocalizeErrors(field.state.meta.errors, t)}
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

            <div>
                <span className="footnote">{t("auth:HAVE_AN_ACCOUNT")}{" "}<Link to="/signin">{t("auth:SIGNIN")}</Link></span>
            </div>
        </form>
    );
};
