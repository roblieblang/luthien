import { useAuth0 } from "@auth0/auth0-react";
import { useUser } from "../../contexts/userContext";

export default function SpotifyAuthButton() {
  const { isAuthenticated } = useAuth0();
  const { userID, spotifyAuthStatus, updateSpotifyAuthStatus } = useUser();

  const handleLogin = () => {
    if (isAuthenticated) {
      fetch("http://localhost:8080/auth/spotify/login")
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
    fetch("http://localhost:8080/auth/spotify/logout", {
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
          className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
          onClick={handleLogout}
        >
          Log Out of Spotify
        </button>
      ) : (
        <button
          className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
          onClick={handleLogin}
        >
          Connect Spotify Account
        </button>
      )}
    </div>
  );
}