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
  const { spotifyAuthStatus } = useUser();

  return (
    <div className="flex flex-col items-center justify-center text-center">
      <BasicHeading text="Home Page" />
      {isAuthenticated && (
        <>
          <div className="flex flex-row space-x-2">
            <Link
              className="rounded-md bg-yellow-400 border-4 border-black hover:bg-yellow-600"
              to="/profile"
            >
              <IoMdPerson size={35} />
            </Link>
            {spotifyAuthStatus && (
              <Link
                className="rounded-md bg-yellow-400 border-4 border-black hover:bg-yellow-600"
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
  );
}
