import { Navbar } from "@/components/features/Navbar";
import { Outlet } from "@tanstack/react-router";

export const NavbarLayout = () => {
    return (
        <>
            <Navbar />
            <Outlet />
        </>
    );
};
