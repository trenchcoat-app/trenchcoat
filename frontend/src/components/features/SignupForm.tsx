import { useMutation } from "@tanstack/react-query";
import { useForm } from "@tanstack/react-form";
import { signUpMutation } from "@/api/@tanstack/react-query.gen";
import { Input, Button } from "@/components/shared";
import type { SignUpBody } from "@/api/types.gen";
import { requiredFieldValidator, confirmPasswordFieldValidator } from "@/utils/validators";
import { composeValidators } from "@/utils/validator-util";

export const SignupForm = () => {
    const mutation = useMutation(signUpMutation());

    const defaultValues: SignUpBody & { confirmPassword: string } = {
        email: "",
        password: "",
        displayName: "",
        confirmPassword: "",
    }

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
                    onBlur: requiredFieldValidator,
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
                        label="Email"
                        placeholder="sample@mail.com"
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>
            
            <form.Field 
                name="displayName"
                validators={{
                    onBlur: requiredFieldValidator,
                }}
            >
                {(field) => (
                    <Input
                        name={field.name}
                        autoComplete={field.name}
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label="Display Name"
                        placeholder="display name"
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Field 
                name="password"
                validators={{
                    onBlur: requiredFieldValidator,
                }}
            >
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

            <form.Field
                name="confirmPassword"
                validators={{
                    onBlur: composeValidators(requiredFieldValidator, confirmPasswordFieldValidator),
                }}
            >
                {(field) => (
                    <Input
                        name={field.name}
                        autoComplete={field.name}
                        type="password"
                        value={field.state.value}
                        onBlur={field.handleBlur}
                        onChange={(e) => field.handleChange(e.target.value)}
                        label="Confirm Password"
                        placeholder="Confirm Password"
                        errors={field.state.meta.errors}
                    />
                )}
            </form.Field>

            <form.Subscribe
                selector={(state) => [state.canSubmit, state.isSubmitting]}
                children={([canSubmit, isSubmitting]) => (
                <Button type="submit" disabled={!canSubmit}>
                    {isSubmitting ? '...' : 'Submit'}
                </Button>
                )}
            />
        </form>
    );
};
