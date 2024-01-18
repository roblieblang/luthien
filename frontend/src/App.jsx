import { useEffect, useState } from "react";
// import { Link, Route, BrowserRouter as Router, Routes } from "react-router-dom";
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

    // if (code) {
    //   fetch(`http://localhost:8080/callback?code=${code}`)
    //   .then(response => response.json())
    //   .then(data => {
    //     // TODO: Handle access token, store it in memory or use it directly to make API calls
    //   });
    // }
  }, []);

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
