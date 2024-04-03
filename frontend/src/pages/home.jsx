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
    <div className="flex flex-col items-center justify-center min-h-screen p-3 -mt-10 text-customButtonText">
      <div className="text-center space-y-4">
        <h1 className="text-xl font-extrabold text-customStroke">
          Seamlessly Sync Your Playlists Between Spotify and YouTube
        </h1>
        <h2 className="text-lg text-customStroke">
          Connect your accounts and let Luthien handle the rest.
        </h2>
      </div>

      <div className="my-8 text-center">
        <h3 className="text-sm font-semibold mb-4">Why Luthien?</h3>
        <ul className="list-none list-inside space-y-2 text-xs">
          <li>Fast and easy playlist conversion</li>
          <li>Secure connection to your music accounts</li>
          <li>Intuitive design for effortless navigation</li>
        </ul>
      </div>

      <div className="flex space-x-3 justify-center items-center my-6">
        <FaSpotify className="text-4xl text-green-600" />
        <FaYoutube className="text-4xl text-red-600" />
        <p className="text-lg">Bring your music together</p>
      </div>

      <div className="absolute inset-0 flex flex-col items-center justify-center -mt-40">
        {isAuthenticated && (
          <>
            {spotifyAuthStatus && youTubeAuthStatus && (
              <Link to="/music" style={{ zIndex: 10 }}>
                <button className="hover:bg-white hover:text-customTertiary transition xl:scale-150 text-sm font-bold rounded bg-customSecondary">
                  <PiSwap size={50} />
                </button>
              </Link>
            )}
            <div className="flex space-x-14">
              <YouTubeAuthButton />
              <SpotifyAuthButton />
            </div>
          </>
        )}
        <div className="xl:scale-150 mt-10">
          <AuthenticationButton />
        </div>
      </div>
      <Footer />
    </div>
  );
}
