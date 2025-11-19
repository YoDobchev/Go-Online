import { defineConfig, loadEnv } from "vite";
import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
    const env = loadEnv(mode, process.cwd(), "");

    const API_BASE = env.VITE_API_URL;
    return {
        plugins: [react()],
        server: {
            proxy: {
                "/api": {
                    target: API_BASE,
                    changeOrigin: true,
                },
            },
        },
    };
});
