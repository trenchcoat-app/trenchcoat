import { createRootRouteWithContext, Outlet } from "@tanstack/react-router";
import type { AuthContextType } from "@/contexts/AuthContext";
import { ToastContainer } from "@/components/features/toast/ToastContainer";
import { NotFound } from "@/components/pages/NotFound";

export interface RouterContext {
    auth: AuthContextType;
}

const RootLayout = () => {
    return (
        <>
            <main style={{ display: "flex", flexDirection: "column", minHeight: "100%", flexGrow: "1" }}>
                <Outlet />
            </main>

            <ToastContainer />
        </>
    );
};

export const rootRoute = createRootRouteWithContext<RouterContext>()({
    component: () => <RootLayout />,
    notFoundComponent: NotFound
});
