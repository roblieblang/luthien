import he from "he";
import { useEffect, useState } from "react";
import { usePlaylist } from "../../contexts/playlistContext";
import LinkButton from "../general/buttons/linkButton";
import TrackList from "../trackList";

export default function SpotifyPlaylist({ playlist }) {
  console.log("spotify playlist: ", playlist.images);
  if (!playlist.images) {
    console.log(`images missing for playlist: ${JSON.stringify(playlist)}`);
  }
  return (
    <div className="bg-slate-600 border border-yellow-600 border-solid p-4 m-2 flex w-10/12">
      <div className="flex-1 text-left">
        {playlist.images && playlist.images[0] && (
          <img
            src={playlist.images[0].url}
            alt={playlist.name}
            className="h-16 w-16 object-cover"
          />
        )}
        <h3 className="text-base font-medium">{he.decode(playlist.name)}</h3>
        <p>Tracks: {playlist.tracks.total}</p>
        <p>Owner: {he.decode(playlist.owner.display_name)}</p>
        <p>{he.decode(playlist.description) || ""}</p>
      </div>
      <div className="ml-auto flex flex-col justify-between">
        <LinkButton
          to={playlist.external_urls.spotify}
          text="Open in Spotify"
        />
        {playlist.tracks.total > 0 && (
          <LinkButton
            to={{
              pathname: "/conversion",
              state: {
                source: "Spotify",
                destination: "YouTube",
                title: playlist.name,
                playlistID: playlist.id,
              },
            }}
            text="Select Playlist"
          />
        )}
      </div>
    </div>
  );
}
