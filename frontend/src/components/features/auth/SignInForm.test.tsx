// SignInForm.test.tsx
import { screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { renderWithQueryClient } from "@/utils/test-util";
import { vi, describe, it, expect } from "vitest";

import { SignInForm } from "./SignInForm";

vi.mock("@tanstack/react-router", () => ({
    useNavigate: () => vi.fn(),
    Link: ({ children }: { children: React.ReactNode }) => <a>{children}</a>,
}));

const setAccount = vi.fn();
vi.mock("@/hooks/useAuth", () => ({
    useAuth: () => ({ setAccount }),
}));

describe("SignInForm", () => {
    it("logs in successfully", async () => {
        const user = userEvent.setup();
        renderWithQueryClient(<SignInForm />);

        await user.type(screen.getByPlaceholderText("auth:EMAIL_PLACEHOLDER"), "test@example.com");
        await user.type(screen.getByPlaceholderText("auth:PASSWORD"), "hunter2");
        await user.click(screen.getByRole("button", { name: "auth:SIGNIN" }));

        await waitFor(() => {
            expect(setAccount).toHaveBeenCalledWith(
                expect.objectContaining({ email: "test@example.com" }),
            );
        });
    });
});