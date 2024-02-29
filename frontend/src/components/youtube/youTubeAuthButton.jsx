import { useAuth0 } from "@auth0/auth0-react";
import { useUser } from "../../contexts/userContext";

export default function YouTubeAuthButton() {
  const { isAuthenticated } = useAuth0();
  const { userID, youTubeAuthStatus, updateYouTubeAuthStatus } = useUser();

  const handleLogin = () => {
    if (isAuthenticated) {
      fetch("http://localhost:8080/auth/google/login")
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
    fetch("http://localhost:8080/auth/google/logout", {
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
          className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
          onClick={handleLogout}
        >
          Log Out of YouTube
        </button>
      ) : (
        <button
          className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
          onClick={handleLogin}
        >
          Connect YouTube Account
        </button>
      )}
    </div>
  );
}
