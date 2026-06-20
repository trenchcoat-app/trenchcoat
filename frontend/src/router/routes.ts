import { createRoute } from "@tanstack/react-router";
import { rootRoute } from "@router/layouts/RootLayout";
import { Home, About, Signup } from "@components/pages";
import { navbarRoute } from "./layouts/NavbarLayout";

export const indexRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/",
    component: Home,
});
export const aboutRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/about",
    component: About,
});
export const signupRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/signup",
    component: Signup
})

export const routeTree = rootRoute.addChildren([
    navbarRoute.addChildren([
        indexRoute, 
        aboutRoute
    ]),
    signupRoute,
]);