import { useEffect, useState } from "react";
import { useUser } from "../contexts/userContext";

const Track = () => {
  /* TODO: Implement for better modularity as styling improves */
};

export default function TrackList({ playlistID }) {
  const [tracks, setTracks] = useState([]);
  const { userID } = useUser();

  const limit = ""; // TODO: will need to handle pagination
  const offset = 0;

  useEffect(() => {
    if (playlistID && userID) {
      fetch(
        // This endpoint seems to be having issues on Spotify's end. Last Working: 2/14/24.
        `http://localhost:8080/spotify/playlist-tracks?userID=${userID}&limit=${limit}&offset=${offset}&playlistID=${playlistID}`
      )
        .then((res) => {
          if (!res.ok) {
            throw new Error("Response from server was not ok");
          }
          return res.json();
        })
        .then((data) => {
          setTracks(data.items);
        })
        .catch((error) => {
          console.error("Error fetching playlist tracks:", error);
        });
    }
  }, [playlistID, userID]);

  if (tracks === null) {
    return <div>Loading tracks...</div>;
  } 

  return (
    <ul>
      {tracks.map((trackItem, index) => (
        <li key={index}>
          <img
            src={trackItem.track.album.images[0].url}
            alt={trackItem.track.name}
            className="h-14 w-14 object-cover"
          />
          {trackItem.track.name} by{" "}
          {trackItem.track.artists.map((artist) => artist.name).join(", ")} from{" "}
          {trackItem.track.album.name}
        </li>
      ))}
    </ul>
  );
}
