import { Auth0Provider } from "@auth0/auth0-react";
import { BrowserRouter as Router } from "react-router-dom";

import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App.jsx";
import { UserProvider } from "./contexts/userContext.jsx";
import "./index.css";

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <Auth0Provider
      domain={import.meta.env.VITE_AUTH0_DOMAIN}
      clientId={import.meta.env.VITE_AUTH0_CLIENT_ID}
      useRefreshTokens={true}
      authorizationParams={{
        scope:
          "openid profile email offline_access https://www.googleapis.com/auth/youtube read:user_idp_tokens",
        redirect_uri: window.location.origin,
      }}
      cacheLocation="localstorage"
    >
      <Router>
        <UserProvider>
          <App />
        </UserProvider>
      </Router>
    </Auth0Provider>
  </React.StrictMode>
);
