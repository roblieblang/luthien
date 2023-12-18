import { useEffect, useState } from "react";
// import {BrowserRouter as Router, Routes, Route, Link } from "react-router-dom"
import "./App.css";

import { SpotifyLoginButton } from "./components/buttons/spotifyLoginButton";
import { SpotifyLogoutButton } from "./components/buttons/spotifyLogoutButton";
import {
  getSpotifyAccessToken,
  getSpotifyRefreshToken,
} from "./utils/spotify-utils";

function App() {
  const [count, setCount] = useState(0);
  const [isAuthenticatedWithSpotify, setIsAuthenticatedWithSpotify] =
    useState(false);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    let code = urlParams.get("code");

    if (isAuthenticatedWithSpotify) {
      getSpotifyRefreshToken();
    }

    if (code) {
      console.log(`Auth code: ${code}`);
      getSpotifyAccessToken(code);
      setIsAuthenticatedWithSpotify(true);
    }
  }, []);

  let accessToken = localStorage.getItem("access_token");
  let refreshToken = localStorage.getItem("refresh_token");
  console.log(`Access: ${accessToken} \n Refresh: ${refreshToken}`);

  return (
    <div className="justify-center text-center">
      <h1 className="text-3xl font-bold underline">Vite + React</h1>
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
        <SpotifyLogoutButton />
      )}
    </div>
  );
}

export default App;
