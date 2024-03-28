import { useEffect } from "react";
import { usePlaylist } from "../../contexts/playlistContext";
import { useUser } from "../../contexts/userContext";
import { config } from "../../utils/config";
import Loading from "./modals/loading";
import { Track } from "./track";

const TrackList = ({ playlistID, sourceType }) => {
  const { tracks, setTracks, setIsFetchingTracks } = usePlaylist();
  const { userID } = useUser();

  useEffect(() => {
    if (playlistID && userID) {
      setIsFetchingTracks(true);
      let url;
      if (sourceType === "spotify") {
        url = `${config.backendUrl}/spotify/playlist-tracks?userID=${userID}&playlistID=${playlistID}`;
      } else if (sourceType === "youtube") {
        url = `${config.backendUrl}/youtube/playlist-tracks?userID=${userID}&playlistID=${playlistID}`;
      }

      if (url) {
        fetch(url)
          .then((res) => {
            if (!res.ok) {
              if (res.status === 401) {
                window.location.href = `/?${sourceType}_session_expired=true`;
                return Promise.reject("Session expired");
              }
              if (sourceType === "youtube" && res.status === 403) {
                window.location.href = `/?${sourceType}_quota_exceeded=true&operation=tracklist-fetch`;
                return Promise.reject("YouTube API quota exceeded");
              }
              throw new Error(
                `Failed to fetch tracks for ${sourceType} playlist ${playlistID}`
              );
            }
            return res.json();
          })
          .then((data) => {
            const formattedTracks = data.items.map((item) => {
              if (sourceType === "spotify") {
                return {
                  id: item.track.uri,
                  title: item.track.name,
                  thumbnailUrl: item.track.album.images[0].url,
                  artist: item.track.artists
                    .map((artist) => artist.name)
                    .join(", "),
                  album: item.track.album.name,
                  link: `https://open.spotify.com/${item.track.uri
                    .split(":")
                    .slice(1)
                    .join("/")}`,
                };
              } else if (sourceType === "youtube") {
                return {
                  id: item.id,
                  title: item.title,
                  thumbnailUrl: item.thumbnailUrl,
                  channelTitle: item.videoOwnerChannelTitle,
                  link: `https://www.youtube.com/watch?v=${item.id}`,
                };
              }
            });
            setTracks(formattedTracks);
          })
          .catch((error) => {
            console.error("Error fetching playlist tracks:", error);
          })
          .finally(() => setIsFetchingTracks(false));
      }
    }
  }, [playlistID, setTracks, userID, sourceType, setIsFetchingTracks]);

  if (!tracks.length) {
    return <Loading />;
  }

  return (
    <div className="flex flex-col items-center px-0">
      {tracks.map((track, index) => (
        <Track key={track.id} track={track} source={sourceType} />
      ))}
    </div>
  );
};

export default TrackList;
