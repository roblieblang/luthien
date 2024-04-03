import { useAuth0 } from "@auth0/auth0-react";
import { useEffect } from "react";
import { FaSpotify, FaYoutube } from "react-icons/fa";
import { PiSwap } from "react-icons/pi";
import { Link, useNavigate } from "react-router-dom";
import "../App.css";
import AuthenticationButton from "../components/auth0/authenticationButton";
import Footer from "../components/general/footer";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyAuthButton from "../components/spotify/spotifyAuthButton";
import YouTubeAuthButton from "../components/youtube/youTubeAuthButton";
import { useUser } from "../contexts/userContext";
import { showErrorToast } from "../utils/toastUtils";

export default function Home() {
  const { isAuthenticated } = useAuth0();
  const { spotifyAuthStatus, youTubeAuthStatus, spotifyUserID } = useUser();
  const navigate = useNavigate();

  useEffect(() => {
    const currentUrl = new URL(window.location.href);
    const queryParams = currentUrl.searchParams;

    const sessionExpiredService = Array.from(queryParams.entries()).find(
      ([key]) => key.endsWith("_session_expired")
    );
    if (sessionExpiredService) {
      const service = sessionExpiredService[0].replace("_session_expired", "");
      showErrorToast(
        `Session expired. Please reauthenticate with ${
          service.charAt(0).toUpperCase() + service.slice(1)
        }`
      );
      setTimeout(() => {
        navigate(window.location.pathname, { replace: true });
        window.location.reload();
      }, 3000);
    }

    const youtubeQuotaExceeded = queryParams.get("youtube_quota_exceeded");
    const operation = queryParams.get("operation");
    if (youtubeQuotaExceeded === "true") {
      let readableOperation = operation || "an operation";
      showErrorToast(
        `YouTube quota exceeded during ${readableOperation}. Please try again later.`
      );
      setTimeout(() => {
        navigate(window.location.pathname, { replace: true });
      }, 3000);
    }
  }, [navigate]);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen p-4 text-customButtonText">
      <div className="text-center space-y-4 max-w-lg mx-auto">
        <h1 className="text-xl sm:text-2xl font-extrabold text-customStroke">
          Seamlessly Sync Your Playlists Between Spotify and YouTube
        </h1>
        <h2 className="text-md sm:text-lg text-customStroke">
          Connect your accounts and let Luthien handle the rest.
        </h2>
      </div>

      <div className="my-8 text-center space-y-4 max-w-lg mx-auto">
        <h3 className="text-sm sm:text-lg font-semibold">Why Luthien?</h3>
        <ul className="list-disc list-inside">
          <li>Fast and easy playlist conversion</li>
          <li>Secure connection to your music accounts</li>
          <li>Intuitive design for effortless navigation</li>
        </ul>
      </div>

      <div className="my-6 text-center space-y-4">
        <div className="flex justify-center items-center text-lg gap-3">
          <FaSpotify className="text-5xl text-green-600" />
          <FaYoutube className="text-5xl text-red-600" />
          <p>Bring your music together</p>
        </div>

        <div className="flex justify-center items-center gap-2 flex-wrap">
          <div className="order-1 sm:order-1">
            <YouTubeAuthButton />
          </div>
          <div className="order-3 sm:order-2">
            <Link
              to="/music"
              className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded inline-flex items-center justify-center"
            >
              <PiSwap className="mr-2" size={24} />
              <span>Start Syncing</span>
            </Link>
          </div>
          <div className="order-2 sm:order-3">
            <SpotifyAuthButton />
          </div>
        </div>
      </div>

      <div className="mt-4 sm:mt-8 mb-10 sm:mb-20">
        <AuthenticationButton />
      </div>

      <Footer />
    </div>
  );
}
