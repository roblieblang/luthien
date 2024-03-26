import he from "he";
import { useState } from "react";
import { usePlaylist } from "../../contexts/playlistContext";
import LinkButton from "../general/buttons/linkButton";
import TrackList from "../trackList";

export default function YouTubePlaylist({ playlist }) {
  const { setPlaylistDetails, tracks, isFetchingTracks } = usePlaylist();
  // TODO: condense UI for both playlist types
  return (
    <div className="bg-slate-600 border border-yellow-600 border-solid p-4 m-2 flex w-10/12">
      <div className="flex-1 text-left">
        {playlist.imageUrl && (
          <img
            src={playlist.imageUrl}
            alt={playlist.title}
            className="h-16 w-16 object-cover"
          />
        )}
        <h3 className="text-base font-medium">{he.decode(playlist.title)}</h3>
        <p>Tracks: {playlist.videosCount}</p>
        {/* TODO: show description on hover like a tooltip */}
        {/* <p>Description: {he.decode(playlist.description) || ""}</p> */}
      </div>
      <div className="ml-auto flex flex-col justify-between">
        <LinkButton
          to={`https://www.youtube.com/playlist?list=${playlist.id}`}
          text="Open on YouTube"
        />
        {playlist.videosCount > 0 && (
          <>
            <LinkButton
              to={{
                pathname: "/conversion",
                state: {
                  source: "YouTube",
                  destination: "Spotify",
                  title: playlist.title,
                  playlistID: playlist.id,
                },
              }}
              text="Select Playlist"
            />
          </>
        )}
      </div>
    </div>
  );
}
