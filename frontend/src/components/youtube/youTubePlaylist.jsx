import { useEffect, useState } from "react";
import { useUser } from "../../contexts/userContext";
import LinkButton from "../general/buttons/linkButton";
import TrackList from "../trackList";

export default function YouTubePlaylist({ playlist }) {
  const [showTracks, setShowTracks] = useState(false);
  const { userID, googleUserID } = useUser();

  const toggleTracks = () => {
    setShowTracks(!showTracks);
  };

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
        <h3 className="text-base font-medium">{playlist.name}</h3>
        <p>{playlist.videosCount} tracks</p>
        {/* TODO: will need to escape certain characters if keeping the description */}
        <p>{playlist.description || ""}</p> 
      </div>
      <div className="ml-auto flex flex-col justify-between">
        <LinkButton
          to={`https://www.youtube.com/playlist?list=${playlist.id}`}
          text="Open in YouTube"
        />
        {/* TODO: these probably shouldn't be LinkButtons */}
        {playlist.videosCount > 0 && (
          <>
            <LinkButton
              to=""
              text={showTracks ? "Hide Tracks" : "View Tracks"}
              onClick={toggleTracks}
            />
            <LinkButton to="" text="Convert" />
          </>
        )}
        {showTracks && <TrackList playlistID={playlist.id} />}
      </div>
    </div>
  );
}
