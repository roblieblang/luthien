import { useEffect } from "react";
import Loading from "../components/general/modals/loading";
import { useUser } from "../contexts/userContext";
import { config } from "../utils/config";

export default function SpotifyCallback() {
  const { userID } = useUser();

  useEffect(() => {
    if (window.location.pathname === "/spotify/callback") {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      // const userID = sessionStorage.getItem("userID"); // I think the retrieval from session storage was faster than context
      const sessionID = sessionStorage.getItem("sessionID");

      fetch(`${config.backendUrl}/auth/spotify/callback`, {
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
  }, [userID]);

  return <Loading />;
}
