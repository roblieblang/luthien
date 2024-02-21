import { useAuth0 } from "@auth0/auth0-react";
import { useEffect, useState } from "react";
import { IoMdPerson } from "react-icons/io";
import { MdLibraryMusic } from "react-icons/md";
import { Link } from "react-router-dom";
import "../App.css";
import AuthenticationButton from "../components/auth0/authenticationButton";
import BasicHeading from "../components/general/headings/basicHeading";
import { LoginButton } from "../components/spotify/loginButton";
import { LogoutButton } from "../components/spotify/logoutButton";
import { useUser } from "../contexts/userContext";

export default function Home() {
  const [isAuthenticatedWithSpotify, setIsAuthenticatedWithSpotify] =
    useState(false);
  const { isAuthenticated } = useAuth0();
  const { userID } = useUser();

  useEffect(() => {
    if (userID) {
      fetch(`http://localhost:8080/auth/spotify/check-auth?userID=${userID}`)
        .then((res) => res.json())
        .then((data) => {
          setIsAuthenticatedWithSpotify(data.isAuthenticated);
        });
    }
  }, [userID]);

  return (
    //This outer div will hold all the components needed in the center of the screen
    <div className="flex flex-col items-center justify-center gap-2 min-h-screen">
      <BasicHeading text="A Playlist Conversion Tool." />
      {isAuthenticated && (
        <>
        
          {/* This div will hold the music button in a row in the center for now */}
          <div className="flex flex-row space-x-2">
            <Link
              className="rounded-md bg-customButton border-4 border-black hover:bg-customStroke"
              to="/profile"
            >
              <IoMdPerson size={45} />
            </Link>
            {isAuthenticatedWithSpotify && (
              <Link
                className="rounded-md bg-customButton border-4 border-black hover:bg-customStroke"
                to="/music"
              >
                <MdLibraryMusic size={45} />
              </Link>
            )}
          </div>

          {/* This div will hold the log in and log out buttons in a row in the center for now */}
          <div className="flex flex-row space-x-2">
            {!isAuthenticatedWithSpotify ? ( // or YouTube
              <LoginButton />
            ) : (
              <LogoutButton
                setIsAuthenticatedWithSpotify={setIsAuthenticatedWithSpotify}
              />
            )}
            <AuthenticationButton />
          </div>
        </>
      )}
    </div>
  );
}
