import { createRoute } from "@tanstack/react-router";
import { rootRoute } from "@/router/routes/RootRoute";
import { navbarRoute } from "@/router/routes/NavbarRoute";
import { protectedRoute } from "@/router/routes/ProtectedRoute";
import { Home, SignUp, SignIn } from "@/components/pages";

export const indexRoute = createRoute({
    getParentRoute: () => navbarRoute,
    path: "/",
    component: Home,
});
export const signUpRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/signup",
    component: SignUp,
});
export const signInRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/signin",
    component: SignIn,
});

export const routeTree = rootRoute.addChildren([
    signUpRoute,
    signInRoute,
    protectedRoute.addChildren([
        navbarRoute.addChildren([
            indexRoute
        ])
    ]),
]);
