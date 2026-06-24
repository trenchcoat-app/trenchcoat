import { signInMutation } from "@/api/@tanstack/react-query.gen";
import type { SignInBody } from "@/api/types.gen";
import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { Input, Button } from "@/components/shared";
import { useTranslation } from "react-i18next";

export const SignInForm = () => {
    const { t } = useTranslation();
    const mutation = useMutation(signInMutation());

    const defaultValues: SignInBody = {
        email: "",
        password: "",
    };

    const form = useForm({
        defaultValues,
        onSubmit: ({ value }) => mutation.mutate({ body: value }),
    });

    return (
        <form
            onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
            }}
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
                        errors={field.state.meta.errors}
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
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
                children={([canSubmit, isSubmitting]) => (
                    <Button type="submit" disabled={!canSubmit}>
                        {isSubmitting ? "..." : t("auth:SIGNIN")}
                    </Button>
                )}
            />
        </form>
    );
};
