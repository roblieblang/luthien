import { useEffect, useState } from "react";
import { FaSpotify, FaYoutube } from "react-icons/fa";
import { useLocation } from "react-router-dom";
import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyUserProfile from "../components/spotify/spotifyUserProfile";
import YouTubeUserProfile from "../components/youtube/youtubeUserProfile";
import { usePlaylist } from "../contexts/playlistContext";
import { showSuccessToast } from "../utils/toastUtils";

export default function Music() {
  const [playlistSource, setPlaylistSource] = useState("");
  const {
    playlistsListCurrentPage,
    setPlaylistsListCurrentPage,
    youTubePlaylistCount,
    spotifyPlaylistCount,
    nextPageToken,
    setNextPageToken,
    prevPageToken,
    setPrevPageToken,
    pageToken,
    setPageToken,
  } = usePlaylist();
  const location = useLocation();

  useEffect(() => {
    const conversionStatus = sessionStorage.getItem("conversionStatus");
    if (conversionStatus && conversionStatus.startsWith("complete:")) {
      const details = conversionStatus.split("complete:")[1];
      showSuccessToast(`Success! ${details}`);
      sessionStorage.removeItem("conversionStatus");
    }
  }, [location]);

  useEffect(() => {
    window.scrollTo({ top: 0, behavior: "smooth" });
  }, [playlistsListCurrentPage]);

  useEffect(() => {
    setPlaylistsListCurrentPage(1);
  }, [playlistSource, setPlaylistsListCurrentPage]);

  const handleNext = () => {
    setPlaylistsListCurrentPage(playlistsListCurrentPage + 1);

    // if (playlistSource === "youtube" && nextPageToken) {
    //   setPageToken(nextPageToken);
    // }
  };

  const handlePrevious = () => {
    if (playlistsListCurrentPage > 1) {
      setPlaylistsListCurrentPage(playlistsListCurrentPage - 1);
    }
    // if (playlistSource === "youtube" && prevPageToken) {
    //   setPageToken(prevPageToken);
    // }
  };

  const handlePageClick = (page) => {
    setPlaylistsListCurrentPage(page);
  };

  const renderPageButtons = () => {
    const buttons = [];
    // First Page Button
    if (playlistsListCurrentPage > 1) {
      buttons.push(
        <button
          key="first"
          onClick={() => handlePageClick(1)}
          className="text-white bg-blue-500 hover:bg-blue-700 font-bold py-1 px-2 rounded"
        >
          1
        </button>
      );
    }
    // Ellipsis for spacing if there are more pages before the two previous pages
    if (playlistsListCurrentPage > 2) {
      buttons.push(
        <span key="ellipsis" className="px-2">
          ...
        </span>
      );
    }
    return buttons;
  };

  const getTotalPages = () => {
    const playlistCount =
      playlistSource === "spotify"
        ? spotifyPlaylistCount
        : youTubePlaylistCount;
    return Math.ceil(playlistCount / 20); // 20 playlists per page
  };

  const totalPages = getTotalPages();

  return (
    <div className="mb-2">
      <div className="mt-2">
        <BasicHeading text="Get Started" />
      </div>
      <div className="-mt-6">
        <h1>Select Playlist Source</h1>
        <div className="flex justify-center p-2">
          <div className="flex items-center space-x-6 text-3xl cursor-pointer">
            <FaSpotify
              className={`lg:scale-125 hover:text-green-400 ${
                playlistSource === "spotify" &&
                "text-green-600 scale-125 lg:scale-150"
              }`}
              aria-label="Spotify"
              role="img"
              onClick={() => setPlaylistSource("spotify")}
            />
            <FaYoutube
              className={`lg:scale-125 hover:text-red-500 ${
                playlistSource === "youtube" &&
                "text-red-600 scale-125 lg:scale-150"
              }`}
              aria-label="YouTube"
              role="img"
              onClick={() => setPlaylistSource("youtube")}
            />
          </div>
        </div>
      </div>
      {playlistSource === "spotify" && <SpotifyUserProfile />}
      {playlistSource === "youtube" && <YouTubeUserProfile />}
      {playlistSource === "spotify" && (
        <div className="flex lg:justify-center lg:space-x-10 justify-between mx-10 mt-2">
          {playlistsListCurrentPage > 1 && (
            <LinkButton text="Previous" onClick={handlePrevious} />
          )}
          {renderPageButtons()}
          <span className="text-lg font-semibold">
            {playlistsListCurrentPage}
          </span>
          {playlistsListCurrentPage < totalPages && (
            <LinkButton text="Next" onClick={handleNext} />
          )}
        </div>
      )}
      <div className="my-5">
        <LinkButton to="/" text="Back" />
      </div>
    </div>
  );
}
