import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { usePlaylist } from "../../../contexts/playlistContext";
import { useUser } from "../../../contexts/userContext";
import { config } from "../../../utils/config";
import { showErrorToast } from "../../../utils/toastUtils";
import LinkButton from "../buttons/linkButton";
import { ModalTrack } from "../track";
import Loading from "./loading";

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
  const navigate = useNavigate();

  const [isLoading, setIsLoading] = useState(false);
  const [adjustedSearchHits, setAdjustedSearchHits] = useState([]);
  const [adjustedSearchMisses, setAdjustedSearchMisses] = useState([]);

  const stockDescription = `Playlist converted from ${source} to ${destination} with Luthien: ${config.frontendUrl}`;

  const fetchSpotifyUserId = async () => {
    fetch(`${config.backendUrl}/spotify/current-profile?userID=${userID}`)
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
        return data.id;
      })
      .catch((error) => {
        console.error("Error fetching Spotify user profile:", error);
      });
  };

  const createNewSpotifyPlaylist = async () => {
    let spotifyUserId;
    if (!spotifyUserID) {
      spotifyUserId = sessionStorage.getItem("spotifyUserId");
      if (!spotifyUserId) {
        spotifyUserId = await fetchSpotifyUserId();
      }
    }
    const url = `${config.backendUrl}/spotify/create-playlist`;
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
      const errorData = await response.json();
      console.error("Error creating playlist:", response.statusText, errorData);
      return;
    }
  };

  const addToNewSpotifyPlaylist = async (newPlaylistId) => {
    const url = `${config.backendUrl}/spotify/add-items-to-playlist`;
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
      const errorData = await response.json();
      console.error(
        "Error adding to playlist:",
        response.statusText,
        errorData
      );
      await rollbackPlaylistCreation(newPlaylistId, "spotify");
      return false;
    }
  };

  const createNewYouTubePlaylist = async () => {
    const url = `${config.backendUrl}/youtube/create-playlist`;
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
      const errorData = await response.json();
      console.error("Error creating playlist:", response.statusText, errorData);
      if (response.status === 403) {
        window.location.href = `/?youtube_quota_exceeded=true&operation=playlist-create`;
      }
      return;
    }
  };

  const addToNewYouTubePlaylist = async (newPlaylistId) => {
    const url = `${config.backendUrl}/youtube/add-items-to-playlist`;
    const videoIds = adjustedSearchHits.map((track) => track.id);
    const payload = {
      userId: userID,
      payload: {
        playlistId: newPlaylistId,
        videoIds: videoIds,
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
      const errorData = await response.json();
      console.error(
        "Error adding to playlist:",
        response.statusText,
        errorData
      );
      await rollbackPlaylistCreation(newPlaylistId, "youtube");
      if (response.status === 403) {
        window.location.href = `/?youtube_quota_exceeded=true&operation=add-to-playlist`;
      }
      return false;
    }
  };

  const rollbackPlaylistCreation = async (newPlaylistId, service) => {
    const url = `${config.backendUrl}/${service}/delete-playlist?userID=${userID}&playlistID=${newPlaylistId}`;
    const response = await fetch(url, {
      method: "DELETE",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (!response.ok) {
      const errorData = await response.json();
      console.error(
        "Error deleting failed conversion:",
        response.statusText,
        errorData
      );
      return false;
    }
    return true;
  };

  const finalizeConversion = async () => {
    setIsLoading(true);
    try {
      let newPlaylistId = await (destination === "Spotify"
        ? createNewSpotifyPlaylist()
        : createNewYouTubePlaylist());

      if (!newPlaylistId) {
        console.error("Failed to create the playlist.");
        showErrorToast(
          "Error during conversion process. Please try again later."
        );
        setIsLoading(false);
        return;
      }

      let tracksAddedSuccessfully = await (destination === "Spotify"
        ? addToNewSpotifyPlaylist(newPlaylistId)
        : addToNewYouTubePlaylist(newPlaylistId));

      if (!tracksAddedSuccessfully) {
        console.error("Failed to add tracks to the playlist.");
        showErrorToast(
          "Error during conversion process. Please try again later."
        );
        setIsLoading(false);
        return;
      }

      // Conversion success
      setPlaylistsLastUpdated(Date.now());
      // Wait for all backend operations to complete before redirecting and refetching playlists
      sessionStorage.setItem(
        "conversionStatus",
        `complete:Converted '${playlistTitle}' from ${source} to ${destination}`
      );
      setTimeout(() => {
        navigate("/music");
      }, 1000);
    } catch (error) {
      showErrorToast("Error during conversion process. Please try again.");
      console.error("Error during the finalize conversion process:", error);
      setIsLoading(false);
    }
  };

  const cancelConversion = () => {
    setAdjustedSearchHits([]);
    setAdjustedSearchMisses([]);
    setTimeout(() => {
      navigate("/music");
    }, 400);
  };

  useEffect(() => {
    if (searchHits && searchHits.length > 0) {
      const transformedHits = searchHits.map((hit) => {
        const newHit = JSON.parse(JSON.stringify(hit[0]));
        return newHit;
      });
      setAdjustedSearchHits(transformedHits);
    }
    if (spotifySearchMisses && spotifySearchMisses.length > 0) {
      const transformedMisses = spotifySearchMisses.map((miss) => {
        const newMiss = JSON.parse(JSON.stringify(miss));
        return newMiss;
      });
      setAdjustedSearchMisses(transformedMisses);
    }
  }, [searchHits, spotifySearchMisses]);

  if (!isOpen) return null;

  if (!adjustedSearchHits || adjustedSearchHits.length === 0 || isLoading) {
    return <Loading />;
  }

  return (
    <div
      className="fixed inset-0 bg-black bg-opacity-70 z-40 flex justify-center items-center"
      onClick={onClose}
    >
      <div
        className="bg-gray-700 p-4 pt-6 rounded-lg lg:w-2/3 w-4/5 max-h-[90vh] overflow-y-auto custom-scrolling-touch"
        style={{ WebkitOverflowScrolling: "touch", overflowY: "auto" }}
        onClick={(e) => e.stopPropagation()}
      >
        <h2 className="text-lg -mt-4 mb-2 text-blue-500 font-bold border-b-2 border-customParagraph">
          {destination} Tracks Found: {adjustedSearchHits.length}
        </h2>
        <div className="flex flex-col items-center">
          {adjustedSearchHits.map((track, index) => (
            <ModalTrack
              key={`${destination}-${track.id}-${index}`}
              track={track}
              destination={destination.toLowerCase()}
              formattedLink={
                destination.toLowerCase() === "youtube"
                  ? `https://www.youtube.com/watch?v=${track.id}`
                  : `https://open.spotify.com/${track.id
                      .split(":")
                      .slice(1)
                      .join("/")}`
              }
              isHit={true}
            />
          ))}
        </div>
        {adjustedSearchMisses.length > 0 && (
          <>
            <h2 className="text-lg mt-1 mb-2 text-red-500 font-bold border-b-2 border-customParagraph">
              Spotify Tracks Not Found: {adjustedSearchMisses.length}
            </h2>
            <div className="flex flex-col items-center">
              {adjustedSearchMisses.map((miss, index) => (
                <ModalTrack
                  track={miss}
                  key={`${miss.songTitle}-${index}`}
                  isHit={false}
                />
              ))}
            </div>
          </>
        )}

        {children}

        <div className="flex flex-col items-center mt-4 space-y-2">
          <LinkButton
            text="Confirm Conversion"
            onClick={finalizeConversion}
            className={
              "hover:bg-white hover:text-green-500 transition text-sm hover:font-extrabold font-bold rounded bg-customSecondary py-1 px-2"
            }
          />
          <LinkButton
            text="Cancel"
            onClick={cancelConversion}
            className={
              "hover:bg-white hover:text-black hover:font-bold transition text-sm font-medium rounded bg-customSecondary py-1 px-2"
            }
          />
        </div>
      </div>
    </div>
  );
}
