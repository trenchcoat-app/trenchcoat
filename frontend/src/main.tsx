import "@/config/i18n";
import "@/config/apiClient";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "@tanstack/react-router";
import { QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider } from "@/contexts/AuthContext";
import { router } from "@/router/router";
import { queryClient } from "@/config/queryClient";
import { useAuth } from "@/hooks/useAuth";
import { ToastProvider } from "@/contexts/ToastContext";
import "@/styles/global.css";

function InnerApp() {
    const auth = useAuth();
    return <RouterProvider router={router} context={{ auth }} />;
}

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <AuthProvider>
                <ToastProvider>
                    <InnerApp />
                </ToastProvider>
            </AuthProvider>
        </QueryClientProvider>
    </StrictMode>,
);
