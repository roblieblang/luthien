import { createContext, useContext, useState } from "react";

const PlaylistContext = createContext();

export const usePlaylist = () => useContext(PlaylistContext);

export const PlaylistProvider = ({ children }) => {
  const [playlistDetails, setPlaylistDetails] = useState(null);
  const [tracks, setTracks] = useState([]);
  const [isFetchingTracks, setIsFetchingTracks] = useState(false);
  const [playlistsLastUpdated, setPlaylistsLastUpdated] = useState(Date.now());
  const [playlistsListCurrentPage, setPlaylistsListCurrentPage] = useState(1);
  const [playlistTracksCurrentPage, setPlaylistTracksCurrentPage] = useState(1);
  const [youTubePlaylistCount, setYouTubePlaylistCount] = useState(0);
  const [spotifyPlaylistCount, setSpotifyPlaylistCount] = useState(0);
  const [nextPageToken, setNextPageToken] = useState("");
  const [prevPageToken, setPrevPageToken] = useState("");
  const [pageToken, setPageToken] = useState("");

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
        playlistsListCurrentPage,
        setPlaylistsListCurrentPage,
        playlistTracksCurrentPage,
        setPlaylistTracksCurrentPage,
        youTubePlaylistCount,
        setYouTubePlaylistCount,
        spotifyPlaylistCount,
        setSpotifyPlaylistCount,
        nextPageToken,
        setNextPageToken,
        prevPageToken,
        setPrevPageToken,
        pageToken,
        setPageToken,
      }}
    >
      {children}
    </PlaylistContext.Provider>
  );
};
