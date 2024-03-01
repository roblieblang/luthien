import he from "he";
import { useState } from "react";
import LinkButton from "../general/buttons/linkButton";
import TrackList from "../trackList";

export default function SpotifyPlaylist({ playlist }) {
  const [showTracks, setShowTracks] = useState(false);

  const toggleTracks = () => {
    setShowTracks(!showTracks);
  };

  return (
    <div className="bg-slate-600 border border-yellow-600 border-solid p-4 m-2 flex w-10/12">
      <div className="flex-1 text-left">
        {playlist.images[0] && (
          <img
            src={playlist.images[0].url}
            alt={playlist.name}
            className="h-16 w-16 object-cover"
          />
        )}
        <h3 className="text-base font-medium">{he.decode(playlist.name)}</h3>
        <p>{playlist.tracks.total} tracks</p>
        <p>Owner: {he.decode(playlist.owner.display_name)}</p>
        <p>{he.decode(playlist.description) || ""}</p>
      </div>
      <div className="ml-auto flex flex-col justify-between">
        <LinkButton
          to={playlist.external_urls.spotify}
          text="Open in Spotify"
        />
        {/* TODO: these probably shouldn't be LinkButtons */}
        {playlist.tracks.total > 0 && (
          <>
            <LinkButton
              to=""
              text={showTracks ? "Hide Tracks" : "View Tracks"}
              onClick={toggleTracks}
            />
            <LinkButton to="" text="Convert" />
          </>
        )}
        {showTracks && (
          <TrackList playlistID={playlist.id} sourceType={"spotify"} />
        )}
      </div>
    </div>
  );
}