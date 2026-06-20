import { createRootRoute, Outlet } from "@tanstack/react-router";

const RootLayout = () => {
    return (
        <main style={{ display: "flex", flexDirection: "column", height: "100%" }}>
            <Outlet />
        </main>
    );
};

export const rootRoute = createRootRoute({
    component: RootLayout,
});
