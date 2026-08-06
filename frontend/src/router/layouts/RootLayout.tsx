import { Outlet } from "@tanstack/react-router";

export const RootLayout = () => {
    return (
        <main style={{ display: "flex", flexDirection: "column", minHeight: "100%" }}>
            <Outlet />
        </main>
    );
};