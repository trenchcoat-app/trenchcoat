import { createRoute } from "@tanstack/react-router";
import { protectedRoute } from "@/router/routes/ProtectedRoute";
import { NavbarLayout } from "@/router/layouts/NavbarLayout";


export const navbarRoute = createRoute({
    getParentRoute: () => protectedRoute,
    id: "navbar-layout",
    component: NavbarLayout,
});
