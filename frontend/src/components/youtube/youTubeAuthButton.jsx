import { useAuth0 } from "@auth0/auth0-react";
import { FaYoutube } from "react-icons/fa";
import { useUser } from "../../contexts/userContext";
import { config } from "../../utils/config";

export default function YouTubeAuthButton() {
  const { isAuthenticated } = useAuth0();
  const { userID, youTubeAuthStatus, updateYouTubeAuthStatus } = useUser();

  const handleLogin = () => {
    if (isAuthenticated) {
      fetch(`${config.backendUrl}/auth/google/login`)
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
    fetch(`${config.backendUrl}/auth/google/logout`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ userID }),
    })
      .then((response) => {
        if (!response.ok) throw new Error("Logout failed");
        updateYouTubeAuthStatus(false);
        window.location.reload(true); // Force the page to refresh upon logout in order to allow for seemless re-auth
      })
      .catch((error) => console.error("Logout error:", error));
  };

  return (
    <div>
      {youTubeAuthStatus ? (
        <button
          className="bg-customHeadline hover:bg-customButton text-md text-customStroke hover:text-red-700 font-bold py-1 px-2 rounded-md border border-black flex items-center justify-center space-x-1"
          onClick={handleLogout}
        >
          <div className="mr-1">
            <FaYoutube aria-label="YouTube" role="img" />
          </div>
          Disconnect
        </button>
      ) : (
        <button
          className="bg-customHeadline hover:bg-customButton text-md text-customStroke hover:text-red-700 font-bold py-1 px-2 rounded-md border border-black flex items-center justify-center space-x-1"
          onClick={handleLogin}
        >
          <div className="mr-1">
            <FaYoutube aria-label="YouTube" role="img" />
          </div>
          Connect
        </button>
      )}
    </div>
  );
}
