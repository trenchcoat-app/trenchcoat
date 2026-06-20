import { useMutation } from "@tanstack/react-query";
import { useForm } from "@tanstack/react-form";
import { signUpMutation } from "@api/@tanstack/react-query.gen";

export const Signup = () => {
    const mutation = useMutation(signUpMutation());

    const form = useForm({
        defaultValues: { name: "", email: "", password: "" },
        onSubmit: ({ value }) => mutation.mutate({ body: value }),
    });

    return (
        <form
            onSubmit={(e) => {
                e.preventDefault();
                form.handleSubmit();
            }}
        >
            <form.Field name="name">
                {(field) => <input name={field.name} value={field.state.value} onBlur={field.handleBlur} onChange={(e) => field.handleChange(e.target.value)} placeholder="Name" />}
            </form.Field>

            <form.Field name="email">
                {(field) => <input name={field.name} type="email" value={field.state.value} onBlur={field.handleBlur} onChange={(e) => field.handleChange(e.target.value)} placeholder="Email" />}
            </form.Field>

            <form.Field name="password">
                {(field) => <input name={field.name} type="password" value={field.state.value} onBlur={field.handleBlur} onChange={(e) => field.handleChange(e.target.value)} placeholder="Password" />}
            </form.Field>

            <button type="submit" disabled={mutation.isPending}>
                {mutation.isPending ? "..." : "Sign Up"}
            </button>
        </form>
    );
};
