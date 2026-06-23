import { Navbar } from "@/components/features/Navbar";
import { createRoute, Outlet } from "@tanstack/react-router";
import { rootRoute } from "@/router/layouts/RootLayout";

const NavbarLayout = () => {
    return (
        <>
            <Navbar />
            <Outlet />
        </>
    );
};

export const navbarRoute = createRoute({
    getParentRoute: () => rootRoute,
    id: "navbar-layout",
    component: NavbarLayout,
});
