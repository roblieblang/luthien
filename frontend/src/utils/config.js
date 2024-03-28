const BACKEND_URL = import.meta.env.VITE_BACKEND_URL || "http://localhost:8080";
const FRONTEND_URL =
  import.meta.env.VITE_FRONTEND_URL || "http://localhost:5173";

export const config = {
  backendUrl: BACKEND_URL,
  frontendUrl: FRONTEND_URL,
};
