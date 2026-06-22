import { createRoute } from "@tanstack/react-router";
import { rootRoute } from "@router/layouts/RootLayout";
import { navbarRoute } from "@router/layouts/NavbarLayout";
import { Home, Signup } from "@components/pages";

export const indexRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/",
    component: Home,
});
export const signupRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/signup",
    component: Signup,
});

export const routeTree = rootRoute.addChildren([navbarRoute.addChildren([indexRoute]), signupRoute]);
