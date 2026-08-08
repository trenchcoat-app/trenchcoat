import { render, screen } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import { Input } from "./Input";

describe("Input", () => {
    it("renders label when provided", () => {
        render(<Input label="Email" />);

        expect(screen.getByText("Email")).toBeInTheDocument();
    });

    it("associates label with input via htmlFor and id", () => {
        render(<Input label="Username" id="username" />);

        const input = screen.getByLabelText("Username");
        expect(input).toBeInTheDocument();
        expect(input).toHaveAttribute("id", "username");
    });

    it("generates an id when none is provided", () => {
        render(<Input label="Password" />);

        const input = screen.getByLabelText("Password");
        expect(input).toBeInTheDocument();
        expect(input.id).toBeTruthy();
    });

    it("renders error message when errors exist", () => {
        render(<Input errors={["Required field"]} />);

        expect(screen.getByRole("alert")).toBeInTheDocument();
        expect(screen.getByText("Required field")).toBeInTheDocument();
    });

    it("sets aria attributes when errors exist", () => {
        render(<Input id="email" errors={["Invalid email"]} />);

        const input = screen.getByRole("textbox");

        expect(input).toHaveAttribute("aria-invalid", "true");
        expect(input).toHaveAttribute("aria-describedby", "email-error");
    });

    it("does not set aria-describedby when no errors exist", () => {
        render(<Input id="email" />);

        const input = screen.getByRole("textbox");

        expect(input).not.toHaveAttribute("aria-describedby");
        expect(input).toHaveAttribute("aria-invalid", "false");
    });

    it("applies invalid class when errors exist", () => {
        render(<Input errors={["Error"]} />);

        const input = screen.getByRole("textbox");

        // css modules hash class names, so check for the original class key
        expect(input.className).toContain("inputInvalid");
    });
});