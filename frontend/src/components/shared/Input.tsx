import type { ComponentProps } from "react";

interface InputProps extends ComponentProps<"input"> {}

export const Input = ({ ref, ...props }: InputProps) => {
    return (
        <input 
            ref={ref} 
            {...props} 
        />
    );
}