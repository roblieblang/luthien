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
            {isAuthenticatedWithSpotify && (
              <Link
                className="rounded-md bg-yellow-400 border-4 border-black hover:bg-yellow-600"
                to="/music"
              >
                <MdLibraryMusic size={35} />
              </Link>
            )}
          </div>
          {!isAuthenticatedWithSpotify ? ( // or YouTube
            <LoginButton />
          ) : (
            <LogoutButton
              setIsAuthenticatedWithSpotify={setIsAuthenticatedWithSpotify}
            />
          )}
        </>
      )}
      <AuthenticationButton />
    </div>
  );
}
