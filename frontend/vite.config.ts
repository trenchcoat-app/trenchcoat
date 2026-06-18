import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vite.dev/config/
export default defineConfig({
    envDir: '../',
    plugins: [react()],
    resolve: {
        alias: {
            "@api": path.resolve(__dirname, "src/api"),
        },
    },
});
