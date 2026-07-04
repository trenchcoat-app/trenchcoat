import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import "@testing-library/jest-dom";
import { Button } from "@/components/shared/Button";

describe("Button", () => {
    it("renders with its children", () => {
        render(<Button>Click me</Button>);

        expect(screen.getByRole("button", { name: "Click me" })).toBeInTheDocument();
    });
});