import { createRoute } from "@tanstack/react-router";
import { rootRoute } from "@/router/layouts/RootLayout";
import { navbarRoute } from "@/router/layouts/NavbarLayout";
import { Home, SignUp, SignIn } from "@/components/pages";

export const indexRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/",
    component: Home,
});
export const signUpRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/signup",
    component: SignUp,
});
export const signInRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/signin",
    component: SignIn,
});

export const routeTree = rootRoute.addChildren([navbarRoute.addChildren([indexRoute, signUpRoute, signInRoute])]);
