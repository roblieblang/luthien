import { useEffect, useState } from "react";
import { usePlaylist } from "../contexts/playlistContext";
import { useUser } from "../contexts/userContext";

const TrackList = ({ playlistID, sourceType }) => {
  const { tracks, setTracks, setIsFetchingTracks } = usePlaylist();
  const { userID } = useUser();

  useEffect(() => {
    if (playlistID && userID) {
      console.log("tracklist pID and uID:", playlistID, userID)
      setIsFetchingTracks(true);
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
              if (res.status === 401) {
                window.location.href = `/?${sourceType}_session_expired=true`;
                return Promise.reject("Session expired");
              }
              throw new Error(
                `Failed to fetch tracks for ${sourceType} playlist ${playlistID}`
              );
            }
            return res.json();
          })
          .then((data) => {
            console.log("tracklist data:", data);
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
                  channelTitle: item.videoOwnerChannelTitle,
                  // artist: "",
                  // album: "",
                };
              }
            });
            setTracks(formattedTracks);
            console.log("tracklist formattedtracks:", formattedTracks);
          })
          .catch((error) => {
            console.error("Error fetching playlist tracks:", error);
          })
          .finally(() => setIsFetchingTracks(false));
      }
    }
  }, [playlistID, setTracks, userID, sourceType, setIsFetchingTracks]);

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
