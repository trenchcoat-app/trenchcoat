import { createRouter, createRootRoute, createRoute, Outlet, Link } from "@tanstack/react-router";
import App from "./App";
import About from "./pages/About";

const rootRoute = createRootRoute({
    component: () => (
        <>
            <nav style={{ display: "flex", gap: "1rem", padding: "1rem", borderBottom: "1px solid #ccc" }}>
                <Link to="/">Home</Link>
                <Link to="/about">About</Link>
            </nav>
            <Outlet />
        </>
    ),
});

const indexRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/",
    component: App,
});

const aboutRoute = createRoute({
    getParentRoute: () => rootRoute,
    path: "/about",
    component: About,
});

const routeTree = rootRoute.addChildren([indexRoute, aboutRoute]);

export const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
    interface Register {
        router: typeof router;
    }
}
