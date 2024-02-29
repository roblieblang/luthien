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
    //This outer div will hold all the components and size according to screen size
    <div className="flex flex-col">

      {/* This header will hold redirect icons in top-left corner of screen */}
      <div className="absolute top-0 left-0 p-4 flex flex-row space-x-2">
        <Link
          className="bg-transparent"
          to="/profile"
        >
          <IoMdPerson size={35} />
        </Link>
        {/* This statement shows the Music Library button if the user is authenticated with Spotify */}
        {isAuthenticatedWithSpotify && (
          <Link
            className="bg-transparent"
            to="/music"
          >
            <MdLibraryMusic size={35} />
          </Link>
        )}
      </div>

      {/* This title will hold the heading component */}
      <title className="absolute inset-0 flex flex-col items-center justify-center mb-10">
        <BasicHeading text="A Playlist Conversion Tool." />
      </title>

       {/* This main content area will hold the login/out buttons */}
      <main className="absolute inset-0 flex flex-col items-center justify-center mt-10">
        {isAuthenticated && (
          <>
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
      </main>

      {/* This content area will be for a minimal footer */}
      <footer>

      </footer>
    </div>
  );
}
