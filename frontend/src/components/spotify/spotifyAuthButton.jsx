import { useAuth0 } from "@auth0/auth0-react";
import { FaSpotify } from "react-icons/fa";
import { useUser } from "../../contexts/userContext";
import { config } from "../../utils/config";

export default function SpotifyAuthButton() {
  const { isAuthenticated } = useAuth0();
  const { userID, spotifyAuthStatus, updateSpotifyAuthStatus } = useUser();

  const handleLogin = () => {
    if (isAuthenticated) {
      fetch(`${config.backendUrl}/auth/spotify/login`)
        .then((response) => {
          if (!response.ok) throw new Error("Network response was not ok");
          return response.json();
        })
        .then((data) => {
          sessionStorage.setItem("sessionID", data.sessionID);
          window.location.href = data.authURL;
        })
        .catch((error) =>
          console.error("There was a problem with the fetch operation:", error)
        );
    }
  };

  const handleLogout = () => {
    fetch(`${config.backendUrl}/auth/spotify/logout`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ userID }),
    })
      .then((response) => {
        if (!response.ok) throw new Error("Logout failed");
        updateSpotifyAuthStatus(false);
        window.location.reload(true); // Force the page to refresh upon logout in order to allow for seemless re-auth
      })
      .catch((error) => console.error("Logout error:", error));
  };

  return (
    <div>
      {spotifyAuthStatus ? (
        <button
          className="bg-customHeadline hover:bg-customButton text-md text-customStroke hover:text-green-600 font-bold py-1 px-2 rounded-md border border-black flex items-center justify-center space-x-1"
          onClick={handleLogout}
        >
          <div className="mr-1">
            <FaSpotify aria-label="YouTube" role="img" />
          </div>
          Disconnect
        </button>
      ) : (
        <button
          className="bg-customHeadline hover:bg-customButton text-md text-customStroke hover:text-green-600 font-bold py-1 px-2 rounded-md border border-black flex items-center justify-center space-x-1"
          onClick={handleLogin}
        >
          <div className="mr-1">
            <FaSpotify aria-label="YouTube" role="img" />
          </div>
          Connect
        </button>
      )}
    </div>
  );
}
