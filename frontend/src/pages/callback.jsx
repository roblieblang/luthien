import { useEffect } from "react";

export default function Callback() {
  useEffect(() => {
    if (window.location.pathname === "/callback") {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const userID = sessionStorage.getItem("userID");
      const sessionID = sessionStorage.getItem("sessionID");

      fetch("http://localhost:8080/auth/spotify/callback", {
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
          console.error("Error during Spotify callback processing:", error);
        });
    }
  }, []);

  return <h1>Loading...</h1>;
}
