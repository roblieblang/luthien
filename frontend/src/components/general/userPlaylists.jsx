import { useEffect, useState } from "react";
import { Bars } from "react-loader-spinner";
import { usePlaylist } from "../../contexts/playlistContext";
import { useUser } from "../../contexts/userContext";
import SpotifyPlaylist from "../spotify/spotifyPlaylist";
import YouTubePlaylist from "../youtube/youTubePlaylist";

const services = {
  spotify: {
    api: ({ userID, offset = 0 }) =>
      `http://localhost:8080/spotify/current-user-playlists?userID=${userID}&offset=${offset}`,
    component: SpotifyPlaylist,
  },
  youtube: {
    api: ({ userID }) =>
      `http://localhost:8080/youtube/current-user-playlists?userID=${userID}`,
    component: YouTubePlaylist,
  },
};

export default function UserPlaylists({ serviceType }) {
  const [playlists, setPlaylists] = useState([]);

  const { userID } = useUser();

  const {
    playlistsLastUpdated,
    playlistsListCurrentPage,
    setPlaylistsListCurrentPage,
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
  } = usePlaylist();

  useEffect(() => {
    const offset = (playlistsListCurrentPage - 1) * 20; // 20 playlists per page

    if (userID && services[serviceType]) {
      const apiURL = services[serviceType].api({
        userID,
        offset,
      });
      fetch(apiURL)
        .then((res) => {
          if (!res.ok) {
            if (res.status === 401) {
              window.location.href = `/?${serviceType}_session_expired=true`;
              return;
            }
            if (serviceType === "youtube" && res.status === 403) {
              window.location.href = `/?${serviceType}_quota_exceeded=true&operation=playlists-fetch`;
              return Promise.reject("YouTube API quota exceeded");
            }
            throw new Error(`Failed to fetch playlists from ${serviceType}`);
          }
          return res.json();
        })
        .then((data) => {
          if (serviceType === "spotify") {
            setSpotifyPlaylistCount(data.total);
          } else {
            setYouTubePlaylistCount(data.totalCount);
            // setNextPageToken(data.nextPageToken || "");
            // setPrevPageToken(data.prevPageToken || "");
          }
          const playlistsData =
            serviceType === "spotify" ? data.items : data.playlists;
          setPlaylists(playlistsData);
        })
        .catch((error) => {
          console.error(`Error fetching ${serviceType} user playlists:`, error);
        });
    }
  }, [
    userID,
    serviceType,
    playlistsLastUpdated,
    playlistsListCurrentPage,
    pageToken,
  ]);

  if (playlists == undefined || playlists.length === 0) {
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

  const PlaylistComponent = services[serviceType]?.component;

  const serviceName = () => {
    if (serviceType.toLowerCase() === "youtube") {
      return "YouTube";
    } else {
      return serviceType.charAt(0).toUpperCase() + serviceType.slice(1);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center">
      <h2 className="text-xl font-semibold">Playlists</h2>
      {playlists.map((playlist) => (
        <PlaylistComponent key={playlist.id} playlist={playlist} />
      ))}
    </div>
  );
}
