import he from "he";
import { useState } from "react";
import { usePlaylist } from "../../contexts/playlistContext";
import LinkButton from "../general/buttons/linkButton";
import TrackList from "../general/trackList";

export default function YouTubePlaylist({ playlist }) {
  return (
    <div className="bg-customBG rounded border-customParagraph border-solid border-2 p-2 my-0.5 flex flex-col w-11/12">
      <div className="flex items-center justify-between mb-2">
        {playlist.imageUrl && (
          <img
            src={playlist.imageUrl}
            alt={playlist.title}
            className="h-16 w-16 object-cover border-2 rounded"
          />
        )}
        <a
          href={`https://www.youtube.com/playlist?list=${playlist.id}`}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block ml-1"
        >
          <h3 className="text-sm font-bold text-slate-200 hover:text-blue-500">
            {he.decode(playlist.title)}
          </h3>
        </a>
        {playlist.videosCount > 0 && (
          <div className="ml-2">
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
              text="Select"
            />
          </div>
        )}
      </div>
      <div className="text-xs ">
        <span>
          {playlist.videosCount}{" "}
          {playlist.videosCount === 1 ? "track" : "tracks"}
        </span>
        <span className="mx-2">•</span>
        <span>Created by {he.decode(playlist.channelTitle)}</span>
        <span className="mx-2">•</span>
        <span>
          {playlist.privacyStatus.charAt(0).toUpperCase() +
            playlist.privacyStatus.slice(1)}
        </span>
      </div>
    </div>
  );
}
