import { useEffect } from "react";
import Loading from "../components/general/modals/loading";
import { config } from "../utils/config";

export default function GoogleCallback() {
  const userID = sessionStorage.getItem("userID");

  useEffect(() => {
    if (window.location.pathname === "/google/callback") {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const sessionID = sessionStorage.getItem("sessionID");

      fetch(`${config.backendUrl}/auth/google/callback`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code, userID, sessionID }),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data.redirectURL) {
            window.location.href = data.redirectURL;
          }
        })
        .catch((error) => {
          console.error("Error during Google callback processing:", error);
        });
    }
  });

  return <Loading />;
}
