import "@/config/i18n";
import "@/config/apiClient";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { QueryClientProvider } from "@tanstack/react-query";
import { AuthProvider } from "@/contexts/AuthProvider";
import { ToastProvider } from "@/contexts/ToastProvider";
import { queryClient } from "@/config/queryClient";
import { App } from "@/App";
import "@/styles/global.css";

async function enableMocking() {
    if (import.meta.env.VITE_ENABLE_MOCKS !== 'true') {
        return;
    }
    
    const { worker } = await import('@/mocks/browser');
    return worker.start();
}

enableMocking().then(() => {
    createRoot(document.getElementById("root")!).render(
        <StrictMode>
            <QueryClientProvider client={queryClient}>
                <AuthProvider>
                    <ToastProvider>
                        <App />
                    </ToastProvider>
                </AuthProvider>
            </QueryClientProvider>
        </StrictMode>,
    );
});
