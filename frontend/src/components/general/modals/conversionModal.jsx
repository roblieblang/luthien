import he from "he";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { usePlaylist } from "../../../contexts/playlistContext";
import { useUser } from "../../../contexts/userContext";
import { showErrorToast } from "../../../utils/toastUtils";
import LinkButton from "../buttons/linkButton";

export default function ConversionModal({
  isOpen,
  onClose,
  children,
  searchHits,
  spotifySearchMisses,
  destination,
  source,
  playlistTitle,
}) {
  const { userID, spotifyUserID, updateSpotifyUserID } = useUser();
  const { setPlaylistsLastUpdated } = usePlaylist();
  const [adjustedSearchHits, setAdjustedSearchHits] = useState([]);
  const [adjustedSearchMisses, setAdjustedSearchMisses] = useState([]);
  const navigate = useNavigate();

  const stockDescription = `Playlist converted from ${source} to ${destination} with Luthien: http://localhost:5173`;

  const fetchSpotifyUserId = async () => {
    fetch(`http://localhost:8080/spotify/current-profile?userID=${userID}`)
      .then((res) => {
        if (!res.ok) {
          if (res.status === 401) {
            window.location.href = `/?spotify_session_expired=true`;
            return;
          }
          throw new Error("Failed to fetch Spotify User Profile and User ID");
        }
        return res.json();
      })
      .then((data) => {
        updateSpotifyUserID(data.id);
        sessionStorage.setItem("spotifyUserId", data.id);
        console.log("fetched spotify user profile id:", data.id);
        return data.id;
      })
      .catch((error) => {
        console.error("Error fetching Spotify user profile:", error);
      });
  };
  // TODO: extract endpoint urls to config file in prod
  const createNewSpotifyPlaylist = async () => {
    let spotifyUserId;
    if (!spotifyUserID) {
      console.log(
        "SpotifyUserID missing from context. Trying session storage..."
      );
      spotifyUserId = sessionStorage.getItem("spotifyUserId");
      if (!spotifyUserId) {
        console.log(
          "spotifyUserId missing from session storage. Fetching a new one..."
        );
        spotifyUserId = await fetchSpotifyUserId();
      }
    }
    const url = `http://localhost:8080/spotify/create-playlist`;
    const payload = {
      userId: userID,
      spotifyUserId: spotifyUserID || spotifyUserId,
      payload: {
        name: playlistTitle,
        public: false,
        collaborative: false,
        description: stockDescription,
      },
    };
    console.log("create spotify playlist payload:", payload);
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (response.ok) {
      const data = await response.json();
      return data.newPlaylistId;
    } else {
      console.error("Error creating playlist:", response.statusText);
      const errorData = await response.json();
      console.error("Error details:", errorData);
      return;
    }
  };

  const addToNewSpotifyPlaylist = async (newPlaylistId) => {
    const url = `http://localhost:8080/spotify/add-items-to-playlist`;
    const trackUris = adjustedSearchHits.map((track) => track.id);
    const payload = {
      userId: userID,
      spotifyPlaylistId: newPlaylistId,
      payload: {
        uris: trackUris,
        position: 0,
      },
    };
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
    if (response.ok) {
      return true;
    } else {
      console.error("Error adding items to playlist:", response.statusText);
      const errorData = await response.json();
      console.error("Error details:", errorData);
      return false;
    }
  };
  // TODO: troubleshoot empty application userid and spotify user id. something's not persisiting in context...
  const createNewYouTubePlaylist = async () => {
    const url = `http://localhost:8080/youtube/create-playlist`;
    const payload = {
      userId: userID,
      payload: {
        title: playlistTitle,
        description: stockDescription,
        privacyStatus: "private",
      },
    };
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });

    if (response.ok) {
      const data = await response.json();
      return data.id;
    } else {
      console.error("Error creating playlist:", response.statusText);
      const errorData = await response.json();
      console.error("Error details:", errorData);
      return;
    }
  };
  // TODO: delete new playlist if addTo<service>Playlist fails
  const addToNewYouTubePlaylist = async (newPlaylistId) => {
    const url = `http://localhost:8080/youtube/add-items-to-playlist`;
    const videoIds = adjustedSearchHits.map((track) => track.id);
    const payload = {
      userId: userID,
      payload: {
        playlistId: newPlaylistId,
        videoIds: videoIds,
      },
    };
    console.log("add items to yt p payload: ", payload);
    const response = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    });
    if (response.ok) {
      return true;
    } else {
      console.error("Error adding items to playlist:", response.statusText);
      const errorData = await response.json();
      console.error("Error details:", errorData);
      return false;
    }
  };

  // TODO: add loading state
  const finalizeConversion = async () => {
    try {
      let newPlaylistId = await (destination === "Spotify"
        ? createNewSpotifyPlaylist()
        : createNewYouTubePlaylist());

      if (!newPlaylistId) {
        console.error("Failed to create the playlist.");
        showErrorToast("Error during conversion process. Please try again.");
        return;
      }

      let tracksAddedSuccessfully = await (destination === "Spotify"
        ? addToNewSpotifyPlaylist(newPlaylistId)
        : addToNewYouTubePlaylist(newPlaylistId));
      if (!tracksAddedSuccessfully) {
        console.error("Failed to add tracks to the playlist.");
        showErrorToast("Error during conversion process. Please try again.");
        return;
      }
      // Conversion success
      setPlaylistsLastUpdated(Date.now());
      // Wait for all backend operations to complete before redirecting and refetching playlists
      sessionStorage.setItem("conversionStatus", "complete");
      setTimeout(() => {
        navigate("/music");
        // TODO: add loading indicator
      }, 1000);
      console.log("Playlist created and tracks added successfully.");
    } catch (error) {
      showErrorToast("Error during conversion process. Please try again.");
      console.error("Error during the finalize conversion process:", error);
    }
  };

  useEffect(() => {
    if (searchHits && searchHits.length > 0) {
      const transformedHits = searchHits.map((hit) => {
        // console.log("hit before transform: ", hit[0]);
        const newHit = JSON.parse(JSON.stringify(hit[0]));
        return newHit;
        // console.log("hit after transform: ", newHit);
      });
      setAdjustedSearchHits(transformedHits);
    }
    if (spotifySearchMisses && spotifySearchMisses.length > 0) {
      console.log("modal misses:", spotifySearchMisses);
      const transformedMisses = spotifySearchMisses.map((miss) => {
        const newMiss = JSON.parse(JSON.stringify(miss));
        return newMiss;
      });
      setAdjustedSearchMisses(transformedMisses);
      console.log("modal adj search misses: ", adjustedSearchMisses);
    }
  }, [searchHits, spotifySearchMisses]);

  if (!isOpen) return null;

  // TODO: display conversion progress

  if (!adjustedSearchHits || adjustedSearchHits.length === 0) {
    return <div>Loading...</div>;
  }

  return (
    <div
      className="fixed inset-0 bg-black bg-opacity-70 z-40 flex justify-center items-center"
      onClick={onClose}
    >
      <div
        className="bg-gray-700 p-6 rounded-lg w-3/4 h-auto sm:h-1/2 lg:max-h-[90vh] overflow-y-auto relative"
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-blue-400">
          Tracks Found: {adjustedSearchHits.length}
        </h2>
        <ul
          className={`max-h-96 overflow-y-auto ${
            adjustedSearchMisses.length > 0 ? "bg-red-900" : ""
          }`}
        >
          {adjustedSearchHits.map((track, index) => (
            <li key={`${destination}-${track.id}-${index}`}>
              <img
                id="track-image"
                src={track.thumbnail}
                alt={he.decode(track.title)}
                className="h-14 w-14 object-cover"
              />
              <div id="track-title">{he.decode(track.title)}</div>
              {track.artist && (
                <div id="track-artist">by {he.decode(track.artist)}</div>
              )}
              {track.album && (
                <div id="track-album">from {he.decode(track.album)}</div>
              )}
            </li>
          ))}
          {adjustedSearchMisses.length > 0 && (
            <div>
              <h2 className="text-blue-400">
                Tracks Not Found on Spotify: {adjustedSearchMisses.length}
              </h2>
              <ul>
                {adjustedSearchMisses.map((miss, index) => (
                  <li key={`miss-${index}`}>
                    {he.decode(miss.songTitle)} by {he.decode(miss.artistName)}
                  </li>
                ))}
              </ul>
            </div>
          )}
        </ul>

        {children}

        <div className="flex flex-col items-center mt-4">
          <LinkButton text="Confirm Conversion" onClick={finalizeConversion} />
        </div>
      </div>
    </div>
  );
}
