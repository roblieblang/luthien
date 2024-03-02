import { useAuth0 } from "@auth0/auth0-react";
import { IoMdPerson } from "react-icons/io";
import { MdLibraryMusic } from "react-icons/md";
import { Link } from "react-router-dom";
import "../App.css";
import AuthenticationButton from "../components/auth0/authenticationButton";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyAuthButton from "../components/spotify/spotifyAuthButton";
import YouTubeAuthButton from "../components/youtube/youTubeAuthButton";
import { useUser } from "../contexts/userContext";

export default function Home() {
  const { isAuthenticated } = useAuth0();
  const { spotifyAuthStatus, youTubeAuthStatus } = useUser();

  return (
    <div className="flex flex-col items-center justify-center text-center">
      <title className="absolute inset-0 flex flex-col items-center justify-center mb-10">
        <BasicHeading text="A Playlist Conversion Tool." />
      </title>
      <div className="absolute inset-0 flex flex-col items-center justify-center mt-10">
        {isAuthenticated && (
          <>
            <div className="absolute top-0 left-0 p-4 flex flex-row space-x-2">
              <Link
                className="bg-transparent"
                to="/profile"
              >
                <IoMdPerson size={35} />
              </Link>
              {(spotifyAuthStatus || youTubeAuthStatus) && (
                <Link
                  className="bg-transparent"
                  to="/music"
                >
                  <MdLibraryMusic size={35} />
                </Link>
              )}
            </div>
            <YouTubeAuthButton />
            <SpotifyAuthButton />
          </>
        )}
        <AuthenticationButton />
      </div>
    </div>
  );
}
