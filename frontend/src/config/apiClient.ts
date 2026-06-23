import { client as apiClient } from "@/api/client.gen";

apiClient.setConfig({
    baseUrl: import.meta.env.VITE_BACKEND_URL || "http://localhost:8080",
});
