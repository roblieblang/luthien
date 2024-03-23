import { createContext, useContext, useState } from "react";

const PlaylistContext = createContext();

export const usePlaylist = () => useContext(PlaylistContext);

export const PlaylistProvider = ({ children }) => {
  const [playlistDetails, setPlaylistDetails] = useState(null);
  const [tracks, setTracks] = useState([]);
  const [isFetchingTracks, setIsFetchingTracks] = useState(false);
  const [playlistsLastUpdated, setPlaylistsLastUpdated] = useState(Date.now());

  const clearTracks = () => setTracks([]);

  return (
    <PlaylistContext.Provider
      value={{
        playlistDetails,
        setPlaylistDetails,
        tracks,
        setTracks,
        clearTracks,
        isFetchingTracks,
        setIsFetchingTracks,
        playlistsLastUpdated,
        setPlaylistsLastUpdated,
      }}
    >
      {children}
    </PlaylistContext.Provider>
  );
};
