import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "@tanstack/react-router";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { client as apiClient } from "./api/client.gen";
import { router } from "./router";
import "./index.css";

// Configure base URL for generated API client from environment variables
apiClient.setConfig({
    baseUrl: import.meta.env.VITE_BACKEND_URL || "http://localhost:8080",
});

const queryClient = new QueryClient();

createRoot(document.getElementById("root")!).render(
    <StrictMode>
        <QueryClientProvider client={queryClient}>
            <RouterProvider router={router} />
        </QueryClientProvider>
    </StrictMode>,
);
