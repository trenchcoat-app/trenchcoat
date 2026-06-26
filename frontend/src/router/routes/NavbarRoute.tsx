import { Navbar } from "@/components/features/Navbar";
import { createRoute, Outlet } from "@tanstack/react-router";
import { protectedRoute } from "@/router/routes/ProtectedRoute";

const NavbarLayout = () => {
    return (
        <>
            <Navbar />
            <Outlet />
        </>
    );
};

export const navbarRoute = createRoute({
    getParentRoute: () => protectedRoute,
    id: "navbar-layout",
    component: NavbarLayout,
});
