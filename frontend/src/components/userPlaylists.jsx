import { useEffect, useState } from "react";
import { useUser } from "../contexts/userContext";
import SpotifyPlaylist from "./spotify/spotifyPlaylist";
import YouTubePlaylist from "./youtube/youTubePlaylist";

const services = {
  spotify: {
    api: ({ userID, limit = "", offset = 0 }) =>
      `http://localhost:8080/spotify/current-user-playlists?userID=${userID}` +
      (limit ? `&limit=${limit}` : "") +
      (offset ? `&offset=${offset}` : ""),
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

  const limit = ""; // TODO: Handle pagination
  const offset = 0;

  useEffect(() => {
    if (userID && services[serviceType]) {
      const { api } = services[serviceType];
      const apiOptions = { userID };

      // Optionally add limit and offset for services that may use them
      if (serviceType === "spotify") {
        apiOptions.limit = limit;
        apiOptions.offset = offset;
      }

      fetch(api(apiOptions))
        .then((res) => {
          if (!res.ok) {
            if (res.status === 401) {
              window.location.href = `/?${serviceType}_session_expired=true`;
              return;
            }
            throw new Error(`Failed to fetch playlists from ${serviceType}`);
          }
          return res.json();
        })
        .then((data) => {
          const playlistsData =
            serviceType === "spotify" ? data.items : data.playlists;
          console.log("YouTube: ", data.playlists);
          console.log("Spotify: ", data.items)
          setPlaylists(playlistsData);
        })
        .catch((error) => {
          console.error(`Error fetching ${serviceType} user playlists:`, error);
        });
    }
  }, [userID, serviceType, limit, offset]);

  if (playlists == undefined || playlists.length === 0) {
    return <div>Loading playlists...</div>;
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
      <h2 className="text-xl font-semibold">{serviceName()} Playlists</h2>
      {playlists.map((playlist) => (
        <PlaylistComponent key={playlist.id} playlist={playlist} />
      ))}
    </div>
  );
}
