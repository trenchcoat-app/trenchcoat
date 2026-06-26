import { createRoute, redirect } from "@tanstack/react-router";
import { rootRoute } from "@/router/routes/RootRoute";

export const protectedRoute = createRoute({
    getParentRoute: () => rootRoute,
    id: "protected",
    beforeLoad: ({ context }) => {
        if (!context.auth.isAuthenticated) {
            throw redirect({ to: "/signin" });
        }
    },
})