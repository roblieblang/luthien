import { useEffect } from "react";
import { useLocation } from "react-router-dom";
import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyUserProfile from "../components/spotify/spotifyUserProfile";
import YouTubeUserProfile from "../components/youtube/youtubeUserProfile";
import { showSuccessToast } from "../utils/toastUtils";

export default function Music() {
  const location = useLocation();

  useEffect(() => {
    if (sessionStorage.getItem("conversionStatus") == "complete") {
      showSuccessToast("Conversion Successful!");
      sessionStorage.removeItem("conversionStatus")
    }
  }, [location]);

  return (
    <div className="mb-2">
      <BasicHeading text="Music Page" />
      <SpotifyUserProfile />
      <YouTubeUserProfile />
      <div className="my-5">
        <LinkButton to="/" text="Back" />
      </div>
    </div>
  );
}
