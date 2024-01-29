import { useAuth0 } from "@auth0/auth0-react";
import { useEffect, useState } from "react";
import { IoMdPerson } from "react-icons/io";
import { Link } from "react-router-dom";
import "../App.css";
import AuthenticationButton from "../components/auth0/authenticationButton";
import { LoginButton } from "../components/spotify/loginButton";
import { LogoutButton } from "../components/spotify/logoutButton";

export default function Home() {
  const [count, setCount] = useState(0);
  const [isAuthenticatedWithSpotify, setIsAuthenticatedWithSpotify] =
    useState(false);
  const { isAuthenticated, user } = useAuth0();

  // TODO: detect when access token has expired and then call on backend refresh token endpoint
  useEffect(() => {
    // Check if user is authenticated with Spotify
    fetch("http://localhost:8080/auth/spotify/check-auth")
      .then((res) => res.json())
      .then((data) => {
        setIsAuthenticatedWithSpotify(data.isAuthenticated);
      });
  }, []);

  useEffect(() => {
    if (user) {
      console.log("Auth0 user ID: ", user.sub);
    }
  }, [user]);

  return (
    <>
      <div className="flex flex-col items-center justify-center text-center">
        <div className="p-2 px-10 bg-blue-600">
          <h1>Homepage</h1>
        </div>
        {isAuthenticated && (
          <Link
            className="rounded-md bg-yellow-400 mt-2 border-4 border-black hover:bg-yellow-600"
            to="/profile"
          >
            <IoMdPerson size={35} />
          </Link>
        )}
        {!isAuthenticatedWithSpotify ? (
          <LoginButton />
        ) : (
          <LogoutButton
            setIsAuthenticatedWithSpotify={setIsAuthenticatedWithSpotify}
          />
        )}
        <AuthenticationButton />
      </div>
    </>
  );
}
