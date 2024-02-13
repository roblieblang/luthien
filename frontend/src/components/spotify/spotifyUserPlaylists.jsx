import { useEffect, useState } from "react";
import SpotifyPlaylist from "./spotifyPlaylist.";

export default function SpotifyUserPlaylists({ userID }) {
  const [playlists, setPlaylists] = useState([]);

  const limit = 5;
  const offset = 0;

  useEffect(() => {
    if (userID) {
      fetch(
        `http://localhost:8080/spotify/current-user-playlists?userID=${userID}&limit=${limit}&offset=${offset}`
      )
        .then((res) => {
          if (!res.ok) {
            throw new Error("Response from server was not ok");
          }
          return res.json();
        })
        .then((data) => {
          setPlaylists(data.items);
        })
        .catch((error) => {
          console.error("Error fetching user playlists:", error);
        });
    }
  }, []);

  if (playlists.length === 0) {
    return <div>Loading playlists...</div>;
  }

  return (
    <div>
      <h2>Playlists</h2>
      {playlists.map((playlist) => (
        <SpotifyPlaylist key={playlist.id} playlist={playlist} />
      ))}
    </div>
  );
}
