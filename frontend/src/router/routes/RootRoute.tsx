import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import type { AuthContextType } from "@/contexts/AuthContext";

export interface RouterContext {
    auth: AuthContextType;
}

const RootLayout = () => {
    return (
        <main style={{ display: "flex", flexDirection: "column", minHeight: "100%" }}>
            <Outlet />
        </main>
    );
};

export const rootRoute = createRootRouteWithContext<RouterContext>()({
    component: () => <RootLayout />,
});
