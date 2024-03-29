import DOMPurify from "dompurify";
import he from "he";
import { useEffect, useState } from "react";
import { usePlaylist } from "../../contexts/playlistContext";
import LinkButton from "../general/buttons/linkButton";
import TrackList from "../general/trackList";

export default function SpotifyPlaylist({ playlist }) {
  return (
    <div className="bg-customBG rounded border-customParagraph border-solid border-2 p-2 my-0.5 flex flex-col lg:w-1/2 w-11/12">
      <div className="flex items-center justify-between mb-2">
        {playlist.images && playlist.images[0] && (
          <img
            src={playlist.images[0].url}
            alt={playlist.name}
            className="lg:h-28 lg:w-28 h-16 w-16 object-cover border-2 rounded"
          />
        )}
        <a
          href={playlist.external_urls.spotify}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block ml-1"
        >
          <h3 className="lg:text-xl text-sm font-bold text-slate-200 hover:text-blue-500">
            {he.decode(playlist.name ?? "")}
          </h3>
        </a>
        {playlist.tracks.total > 0 && (
          <div className="ml-2">
            <LinkButton
              to={{
                pathname: "/conversion",
                state: {
                  source: "Spotify",
                  destination: "YouTube",
                  title: playlist.name,
                  playlistID: playlist.id,
                  size: playlist.tracks.total,
                },
              }}
              text="Select"
            />
          </div>
        )}
      </div>

      {/* Bottom Row: Track Total and Owner */}
      <div className="xl:-mt-2 lg:text-lg text-xs text-white">
        <span>
          {playlist.tracks.total}{" "}
          {playlist.tracks.total === 1 ? "track" : "tracks"}
        </span>
        <span className="mx-2">•</span>
        <span>Created by {he.decode(playlist.owner.display_name ?? "")}</span>
        <span className="mx-2">•</span>
        <span>{playlist.public ? "Public" : "Private"}</span>
      </div>
    </div>
  );
}
