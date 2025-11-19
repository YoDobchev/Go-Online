import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { API_BASE } from "./src/config"

// https://vite.dev/config/
export default defineConfig({
    plugins: [react()],
    server: {
        proxy: {
            "/api": {
                target: API_BASE,
                changeOrigin: true,
            },
        },
    },
});
