import he from "he";
import { useEffect, useState } from "react";
import { Bars } from "react-loader-spinner";
import { usePlaylist } from "../../contexts/playlistContext";
import { useUser } from "../../contexts/userContext";

const Track = ({ track, source }) => {
  return (
    <div className="bg-customBG rounded border-customParagraph border-solid border-2 p-2 my-0.5 flex w-11/12">
      <div className="flex-none">
        {track.thumbnailUrl && (
          <img
            src={track.thumbnailUrl}
            alt={track.title}
            className="lg:h-28 lg:w-28 h-14 w-14 object-cover mr-2 border-2 rounded"
          />
        )}
      </div>
      <div className="flex-1 text-center flex flex-col justify-center">
        <a
          href={track.link}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block"
        >
          <h3 className="lg:text-2xl text-sm font-bold text-slate-200 hover:text-blue-500">
            {he.decode(track.title)}
          </h3>
        </a>
        <div className="lg:text-lg text-xs">
          <span>{source.charAt(0).toUpperCase() + source.slice(1)}</span>
          <span className="mx-2">•</span>
          {source === "spotify" ? (
            <>
              <span>{track.artist}</span>
              <span className="mx-2">•</span>
              <span>{track.album}</span>{" "}
            </>
          ) : (
            <span>{track.channelTitle}</span>
          )}
        </div>
      </div>
    </div>
  );
};

const TrackList = ({ playlistID, sourceType }) => {
  const { tracks, setTracks, setIsFetchingTracks } = usePlaylist();
  const { userID } = useUser();

  useEffect(() => {
    if (playlistID && userID) {
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
    return (
      <div className="flex h-screen items-center justify-center">
        <Bars
          height="80"
          width="80"
          color="#e2714a"
          ariaLabel="bars-loading"
          visible={true}
        />
      </div>
    );
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
