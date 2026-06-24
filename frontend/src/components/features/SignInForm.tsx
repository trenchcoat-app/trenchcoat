import { signInMutation } from "@/api/@tanstack/react-query.gen";
import type { SignInBody } from "@/api/types.gen";
import { useForm } from "@tanstack/react-form";
import { useMutation } from "@tanstack/react-query";
import { Input, Button } from "@/components/shared";

export const SignInForm = () => {
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
                        label="Email"
                        placeholder="sample@mail.com"
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Field name="password">
                {(field) => (
                    <Input
                        name={field.name}
                        autoComplete={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label="Password"
                        placeholder="Password"
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
                children={([canSubmit, isSubmitting]) => (
                    <Button type="submit" disabled={!canSubmit}>
                        {isSubmitting ? "..." : "Submit"}
                    </Button>
                )}
            />
        </form>
    );
};
