import { RouterProvider } from "@tanstack/react-router";
import { router } from "@/router/router";
import { useAuth } from "@/hooks/useAuth";

export const App = () => {
    const auth = useAuth();

    return <RouterProvider router={router} context={{ auth }} />;
}