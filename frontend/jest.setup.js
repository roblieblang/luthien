/* eslint-disable no-undef */
require("dotenv").config({ path: "./.env.test" });

Object.defineProperty(global, "importMeta", {
  value: {
    env: {
      VITE_BACKEND_URL: process.env.VITE_BACKEND_URL,
      VITE_PORT: process.env.VITE_PORT,
      VITE_AUTH0_DOMAIN: process.env.VITE_AUTH0_DOMAIN,
      VITE_AUTH0_CLIENT_ID: process.env.VITE_AUTH0_CLIENT_ID,
    },
  },
});
