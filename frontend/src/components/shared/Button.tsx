import type { ComponentProps } from "react";
import "./Button.css";

interface ButtonProps extends ComponentProps<"button"> {}

export const Button = ({ ref, ...props }: ButtonProps) => {
    return (
        <div className="button-wrapper">
            <button className="button" ref={ref} {...props} />
        </div>
    );
};
