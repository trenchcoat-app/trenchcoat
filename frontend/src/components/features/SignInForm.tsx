import { signInMutation } from "@/api/@tanstack/react-query.gen";
import type { SignInBody, SignInOkResponse } from "@/api/types.gen";
import { useForm } from "@tanstack/react-form";
import { Link, useNavigate } from "@tanstack/react-router";
import { useMutation } from "@tanstack/react-query";
import { Input, Button, AuthFormLayout } from "@/components/shared";
import { useTranslation } from "react-i18next";
import { useAuth } from "@/hooks/useAuth";

export const SignInForm = () => {
    const { t } = useTranslation();
    const navigate = useNavigate();
    const { setAccount } = useAuth();

    const mutation = useMutation({
        ...signInMutation(),
        onSuccess: (data: SignInOkResponse) => {
            setAccount(data.account);
            navigate({ to: "/" });
        },
    });

    const defaultValues: SignInBody = {
        email: "",
        password: "",
    };

    const form = useForm({
        defaultValues,
        onSubmit: ({ value }) => mutation.mutate({ body: value }),
    });

    return (
        <AuthFormLayout
            title={t("auth:SIGNIN_TITLE")}
            onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
            }}
            note={
                <span>
                    {t("auth:DONT_HAVE_AN_ACCOUNT")} <Link to="/signup">{t("auth:SIGNUP")}</Link>
                </span>
            }
        >
            <form.Field name="email">
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
                    />
                )}
            </form.Field>

            <form.Field name="password">
                {(field) => (
                    <Input
                        name={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label={t("auth:PASSWORD")}
                        placeholder={t("auth:PASSWORD")}
                    />
                )}
            </form.Field>

            <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
                children={([canSubmit, isSubmitting]) => (
                    <Button type="submit" disabled={!canSubmit || isSubmitting}>
                        {t("auth:SIGNIN")}
                    </Button>
                )}
            />
        </AuthFormLayout>
    );
};
