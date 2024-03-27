import { useEffect } from "react";
import { Bars } from "react-loader-spinner";
import { useUser } from "../contexts/userContext";

export default function SpotifyCallback() {
  const { userID } = useUser();

  useEffect(() => {
    if (window.location.pathname === "/spotify/callback") {
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      // const userID = sessionStorage.getItem("userID"); // I think the retrieval from session storage was faster than context
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
  }, [userID]);

  return (
    <div className="flex h-screen items-center justify-center">
      <Bars
        height="80"
        width="80"
        color="#e2714a"
        ariaLabel="bars-loading"
        visible={true}
      />
    </div>
  );
}
