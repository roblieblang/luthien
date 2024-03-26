import { useEffect, useRef, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";
import ConversionModal from "../components/general/modals/conversionModal";
import TrackList from "../components/trackList";
import { usePlaylist } from "../contexts/playlistContext";
import { useUser } from "../contexts/userContext";

export default function Conversion() {
  const [conversionProgress, setConversionProgress] = useState(0);
  const [conversionStatus, setConversionStatus] = useState("searching");
  const [spotifySearchMisses, setSpotifySearchMisses] = useState([]); // There are no search misses for youtube

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchHits, setSearchHits] = useState([]);
  const { tracks, clearTracks } = usePlaylist();
  const location = useLocation();
  const navigate = useNavigate();
  const { userID } = useUser();
  const { source, destination, title, playlistID } = location.state || {};

  // TODO: manage loading states with `isFetchingTracks`

  const constructSpotifySearchUrlsUsingVideoTitles = (videoTitles) => {
    return videoTitles.map(
      (title) =>
        `http://localhost:8080/spotify/search-using-video?userID=${userID}&videoTitle=${title}`
    );
  };

  const constructYouTubeSearchUrls = () => {
    return tracks.map(
      (track) =>
        `http://localhost:8080/youtube/search-for-video?userID=${userID}&songTitle=${track.title}&artistName=${track.artist}`
    );
  };

  const initiateSearch = async (searchUrls, artistSongPairs) => {
    try {
      const results = await Promise.all(
        searchUrls.map(async (url, index) => {
          try {
            const res = await fetch(url);
            if (res.status === 401) {
              navigate(`/?${destination}_session_expired=true`, {
                replace: true,
              });
              throw new Error("401 Unauthorized");
            }
            if (res.status === 404) {
              setSpotifySearchMisses((prevMisses) => [
                ...prevMisses,
                artistSongPairs[index],
              ]);
            }
            if (!res.ok) throw new Error(`Failed to fetch: ${res.statusText}`);
            const data = await res.json();
            if (data.error && data.error === "No tracks found") {
              setSpotifySearchMisses((prevMisses) => [
                ...prevMisses,
                artistSongPairs[index],
              ]);
              return null;
            }
            return data;
          } catch (error) {
            console.error(error);
            return null;
          }
        })
      );

      // Filter out null values (misses) and process search hits
      const hits = results.filter((result) => result !== null);
      setSearchHits((prevHits) => [...prevHits, ...hits]);
    } catch (error) {
      console.error("Error searching for tracks", error);
    }
  };

  const handleConvertClick = () => {
    setIsModalOpen(true);
    if (destination === "YouTube") {
      const searchUrls = constructYouTubeSearchUrls();

      const artistSongPairs = tracks.map((track) => [
        track.artist,
        track.title,
      ]);
      initiateSearch(searchUrls, artistSongPairs); // Adjusted to pass artistSongPairs
    } else {
      const videoTitles = tracks.map((track) => track.title);
      const cleanedVideoTitles = videoTitles.map((title) =>
        // eslint-disable-next-line no-useless-escape
        title.replace(/[.,\/#!$%\^&\*;:{}=\-_`'~()\[\]【】『』]/g, "").trim()
      );

      for (let i = 0; i < videoTitles.length; i++) {
        console.log(
          `before cleaning: ${videoTitles[i]} => after: ${cleanedVideoTitles[i]}`
        );
      }
      const searchUrls =
        constructSpotifySearchUrlsUsingVideoTitles(cleanedVideoTitles);
      const pairs = cleanedVideoTitles.map((title) => ({
        songTitle: title,
        artistName: "",
      }));
      initiateSearch(searchUrls, pairs);
    }
  };

  useEffect(() => {
    return () => {
      clearTracks();
    };
  }, [location.pathname]);

  useEffect(() => {
    let missingValues = [];
    if (!source) missingValues.push("source");
    if (!title) missingValues.push("title");
    if (!destination) missingValues.push("destination");

    if (missingValues.length > 0) {
      console.log(`Missing value(s): ${missingValues.join(", ")}`);
      navigate("/music");
    }
  }, [source, title, destination, navigate]);

  return (
    <div className="py-5 mb-5">
      <BasicHeading
        text={`Convert'${title}' from ${source} to ${destination}`}
      />
      <TrackList playlistID={playlistID} sourceType={source.toLowerCase()} />
      <LinkButton text="Convert" onClick={handleConvertClick} />
      <div className="my-5">
        <LinkButton to="/music" text="Back" />
      </div>
      <ConversionModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        searchHits={searchHits}
        spotifySearchMisses={spotifySearchMisses}
        source={source}
        destination={destination}
        playlistTitle={title}
      ></ConversionModal>
    </div>
  );
}
