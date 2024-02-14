import { useEffect, useState } from "react";
import { useUser } from "../../contexts/userContext";
import SpotifyPlaylist from "./spotifyPlaylist.";

export default function SpotifyUserPlaylists() {
  const [playlists, setPlaylists] = useState([]);
  const { userID } = useUser();

  const limit = ""; // TODO: will need to handle pagination
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
  }, [userID]);

  if (playlists.length === 0) {
    return <div>Loading playlists...</div>;
  }

  return (
    <div className="flex flex-col items-center justify-center">
      <h2 className="text-xl font-semibold">Playlists</h2>
      {playlists.map((playlist) => (
        <SpotifyPlaylist key={playlist.id} playlist={playlist} />
      ))}
    </div>
  );
}
