import { useEffect, useState } from "react";
import { useUser } from "../contexts/userContext";

const TrackList = ({ playlistID, sourceType }) => {
  const [tracks, setTracks] = useState([]);
  const { userID } = useUser();

  useEffect(() => {
    if (playlistID && userID) {
      let url;
      if (sourceType === "spotify") {
        url = `http://localhost:8080/spotify/playlist-tracks?userID=${userID}&playlistID=${playlistID}`;
      } else if (sourceType === "youtube") {
        url = `http://localhost:8080/youtube/playlist-tracks?userID=${userID}&playlistID=${playlistID}`;
      }

      if (url) {
        fetch(url)
          .then((res) => {
            if (!res.ok) {
              throw new Error("Response from server was not ok");
            }
            return res.json();
          })
          .then((data) => {
            const formattedTracks = data.items.map((item) => {
              if (sourceType === "spotify") {
                return {
                  id: item.track.id,
                  title: item.track.name,
                  thumbnailUrl: item.track.album.images[0].url,
                  artist: item.track.artists
                    .map((artist) => artist.name)
                    .join(", "),
                  album: item.track.album.name,
                };
              } else if (sourceType === "youtube") {
                return {
                  id: item.id,
                  title: item.title,
                  thumbnailUrl: item.thumbnailUrl,
                  artist: "",
                  album: "",
                };
              }
            });
            console.log(formattedTracks);
            setTracks(formattedTracks);
          })
          .catch((error) => {
            console.error("Error fetching playlist tracks:", error);
          });
      }
    }
  }, [playlistID, userID, sourceType]);

  if (!tracks.length) {
    return <div>Loading tracks...</div>;
  }

  return (
    <ul>
      {tracks.map((track, index) => (
        <li key={index}>
          <img
            src={track.thumbnailUrl}
            alt={track.title}
            className="h-14 w-14 object-cover"
          />
          <div>{track.title}</div>
          {track.artist && <div>by {track.artist}</div>}
          {track.album && <div>from {track.album}</div>}
        </li>
      ))}
    </ul>
  );
};

export default TrackList;
