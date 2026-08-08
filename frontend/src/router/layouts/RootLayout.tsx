import { Outlet } from "@tanstack/react-router";
import { ToastContainer } from "@/components/features/toast/ToastContainer";

export const RootLayout = () => {
    return (
        <>
            <main style={{ display: "flex", flexDirection: "column", minHeight: "100%", flexGrow: "1" }}>
                <Outlet />
            </main>

            <ToastContainer />
        </>
    );
};