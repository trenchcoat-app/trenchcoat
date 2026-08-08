import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Button } from "@/components/shared/Button";

describe("Button", () => {
    it("renders with its children", () => {
        render(<Button>Click me</Button>);

        expect(screen.getByRole("button", { name: "Click me" })).toBeInTheDocument();
    });

    it("calls onClick when clicked", async () => {
        const user = userEvent.setup();
        const onClick = vi.fn();
        render(<Button onClick={onClick}>Click me</Button>);

        await user.click(screen.getByRole("button", { name: "Click me" }));

        expect(onClick).toHaveBeenCalledTimes(1);
    });

    it("does not call onClick when aria-disabled is set", async () => {
        const user = userEvent.setup();
        const onClick = vi.fn();
        render(<Button aria-disabled onClick={onClick}>Click me</Button>);

        await user.click(screen.getByRole("button", { name: "Click me" }));

        expect(onClick).not.toHaveBeenCalled();
    });

    it("passes native button props", () => {
        render(<Button disabled>Click me</Button>);

        expect(screen.getByRole("button", { name: "Click me" })).toBeDisabled();
    });

    it("passes type prop", () => {
        render(<Button type="submit">Click me</Button>);

        expect(screen.getByRole("button", { name: "Click me" })).toHaveAttribute("type", "submit");
    });

    it("passes aria attributes", () => {
        render(<Button aria-label="Custom label">Click me</Button>);

        expect(screen.getByRole("button", { name: "Custom label" })).toBeInTheDocument();
    });

    it("forwards ref to the button element", () => {
        const ref = { current: null };
        render(<Button ref={ref}>Click me</Button>);

        expect(ref.current).toBeInstanceOf(HTMLButtonElement);
    });
});