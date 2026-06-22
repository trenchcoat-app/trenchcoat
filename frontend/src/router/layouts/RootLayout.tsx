import { createRootRoute, Outlet } from "@tanstack/react-router";

const RootLayout = () => {
    return (
        <main>
            <Outlet />
        </main>
    );
};

export const rootRoute = createRootRoute({
    component: RootLayout,
});
