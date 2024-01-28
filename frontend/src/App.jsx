import { useEffect, useState } from "react";
import "./App.css";

import { SpotifyLoginButton } from "./components/buttons/spotifyLoginButton";
import { SpotifyLogoutButton } from "./components/buttons/spotifyLogoutButton";

function App() {
  const [count, setCount] = useState(0);
  const [isAuthenticatedWithSpotify, setIsAuthenticatedWithSpotify] =
    useState(false);

  // TODO: detect when access token has expired and then call on refresh token endpoint
  
  useEffect(() => {
    // Check if user is authenticated with Spotify
    fetch("http://localhost:8080/auth/spotify/check-auth")
      .then((res) => res.json())
      .then((data) => {
        setIsAuthenticatedWithSpotify(data.isAuthenticated);
      });
  }, []);

  return (
    <div className="justify-center text-center">
      <h1 className="text-3xl font-bold underline">Vite + react</h1>
      <div className="card">
        <button
          onClick={() => {
            setCount((count) => count + 1);
            alert("Eyoy sire");
          }}
        >
          count is {count}
        </button>
      </div>
      {!isAuthenticatedWithSpotify ? (
        <SpotifyLoginButton />
      ) : (
        <SpotifyLogoutButton
          setIsAuthenticatedWithSpotify={setIsAuthenticatedWithSpotify}
        />
      )}
    </div>
  );
}

export default App;
