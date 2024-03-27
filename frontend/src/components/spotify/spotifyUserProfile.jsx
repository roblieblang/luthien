import { useEffect, useState } from "react";
import { Bars } from "react-loader-spinner";
import { useUser } from "../../contexts/userContext";
import UserPlaylists from "../general/userPlaylists";

export default function SpotifyUserProfile() {
  const { userID, updateSpotifyUserID } = useUser();
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    if (userID) {
      fetch(`http://localhost:8080/spotify/current-profile?userID=${userID}`)
        .then((res) => {
          if (!res.ok) {
            if (res.status === 401) {
              window.location.href = `/?spotify_session_expired=true`;
              return;
            }
            throw new Error("Failed to fetch Spotify User Profile");
          }
          return res.json();
        })
        .then((data) => {
          setProfile(data);
          updateSpotifyUserID(data.id);
          sessionStorage.setItem("spotifyUserId", data.id);
        })
        .catch((error) => {
          console.error("Error fetching user profile:", error);
        });
    }
  }, [userID]);

  if (!profile) {
    return (
      <div className="flex h-screen items-center justify-center">
        <Bars
          height="80"
          width="80"
          color="#e2714a"
          ariaLabel="bars-loading"
          visible={true}
        />
      </div>
    );
  }

  return (
    <div>
      <h1 className="font-bold text-lg">
        {profile.display_name}&apos;s Spotify Account
      </h1>

      <div className="mt-2">
        <UserPlaylists serviceType={"spotify"} />
      </div>
    </div>
  );
}
