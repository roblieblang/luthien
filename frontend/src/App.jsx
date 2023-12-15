import { useEffect, useState } from "react";
// import {BrowserRouter as Router, Routes, Route, Link } from "react-router-dom"
import "./App.css";

import { SpotifyLoginButton } from "./components/spotifyLoginButton";

function App() {
  const [count, setCount] = useState(0);
  const [isAuthenticatedWithSpotify, setIsAuthenticatedWithSpotify] =
    useState(false);

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    let code = urlParams.get("code");
    console.log(`Spotify auth token: ${code}`);

    if (code) {
      setIsAuthenticatedWithSpotify(!isAuthenticatedWithSpotify);
    }
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
        <h1 className="text-green-500">You are authenticated with Spotify.</h1>
      )}
    </div>
  );
}

export default App;
