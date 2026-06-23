import type { ComponentProps } from "react";

interface ButtonProps extends ComponentProps<"button"> {}

export const Button = ({ ref, ...props }: ButtonProps) => {
    return <button ref={ref} {...props} />;
};
