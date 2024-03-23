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

  const fetchInitiatedRef = useRef(false);

  // TODO: manage loading states with `isFetchingTracks`
  // TODO: handle 3rd party usage quotas on frontend (as well as backend)

  const fetchArtistAndSong = (videoTitles) => {
    return fetch("http://localhost:8080/auth/openai/extract-artist-song", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ videoTitles }),
    })
      .then((response) => response.json())
      .catch((error) => console.error("Error:", error));
  };

  const constructSpotifySearchUrls = (response, userID) => {
    return response.result.map(({ artistName, songTitle }) => {
      let url = `http://localhost:8080/spotify/search-for-track?userID=${userID}&limit=1`;

      if (songTitle) {
        url += `&trackTitle=${songTitle}`;
      }

      if (artistName) {
        url += `&artistName=${artistName}`;
      }

      return url;
    });
  };

  useEffect(() => {
    if (!tracks.length || fetchInitiatedRef.current) return;

    fetchInitiatedRef.current = true;
    console.log("Trying to search for tracks...");

    const constructYouTubeSearchUrls = () =>
      tracks.map(
        (track) =>
          `http://localhost:8080/youtube/search-for-video?userID=${userID}&songTitle=${track.title}&artistName=${track.artist}`
      );

    const initiateSearch = async (searchUrls, artistSongPairs) => {
      console.log("searchUrls:", searchUrls);
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
              if (!res.ok)
                throw new Error(`Failed to fetch: ${res.statusText}`);
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

    if (destination === "YouTube") {
      const searchUrls = tracks.map(
        (track) =>
          `http://localhost:8080/youtube/search-for-video?userID=${userID}&songTitle=${track.title}&artistName=${track.artist}`
      );

      const artistSongPairs = tracks.map((track) => [
        track.artist,
        track.title,
      ]);
      initiateSearch(searchUrls, artistSongPairs); // Adjusted to pass artistSongPairs
    } else {
      const videoTitles = tracks.map((track) => track.title);
      fetchArtistAndSong(videoTitles)
        .then((response) => {
          const searchUrls = constructSpotifySearchUrls(response, userID);
          initiateSearch(searchUrls, response.result);
        })
        .catch((error) =>
          console.error("Error fetching artist and song:", error)
        );
    }
  }, [destination, title, userID, tracks]);

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

  const handleOpenModal = () => {
    setIsModalOpen(true);
  };
  
  console.log("tracks:", tracks);

  console.log("searchHits", searchHits);
  console.log("conversion page search misses:", spotifySearchMisses);

  return (
    <div className="py-5 mb-5">
      <BasicHeading
        text={`Convert'${title}' from ${source} to ${destination}`}
      />
      <TrackList playlistID={playlistID} sourceType={source.toLowerCase()} />
      <LinkButton text="Convert" onClick={handleOpenModal} />
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
