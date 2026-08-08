import { createRootRouteWithContext } from "@tanstack/react-router";
import type { AuthContextType } from "@/contexts/AuthContext";
import { RootLayout } from "@/router/layouts/RootLayout";
import { NotFound } from "@/components/pages/NotFound";

export interface RouterContext {
    auth: AuthContextType;
}

export const rootRoute = createRootRouteWithContext<RouterContext>()({
    component: () => <RootLayout />,
    notFoundComponent: NotFound,
});
