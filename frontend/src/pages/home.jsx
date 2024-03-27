import { useAuth0 } from "@auth0/auth0-react";
import { useEffect } from "react";
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
    <div className="flex items-center justify-center text-center">
      <div className="mt-5">
        <BasicHeading text="Convert Your Playlists" />
        <p className="text-xs -mt-10">(Demo App)</p>
      </div>
      <div className="absolute inset-0 flex flex-col items-center justify-center -mt-40">
        {isAuthenticated && (
          <>
            {spotifyAuthStatus && youTubeAuthStatus && (
              <Link to="/music" style={{ zIndex: 10 }}>
                <button className="hover:bg-white hover:text-customTertiary transition text-sm font-bold rounded bg-customSecondary">
                  <PiSwap size={50} />
                </button>
              </Link>
            )}
            <div className="flex space-x-20 -my-11">
              <YouTubeAuthButton />
              <SpotifyAuthButton />
            </div>
          </>
        )}
        <div className="mt-20">
          <AuthenticationButton />
        </div>
      </div>
      <Footer />
    </div>
  );
}
