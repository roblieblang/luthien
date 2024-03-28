import { useEffect, useState } from "react";
import { useUser } from "../../contexts/userContext";
import { config } from "../../utils/config";
import Loading from "../general/modals/loading";
import UserPlaylists from "../general/userPlaylists";

export default function SpotifyUserProfile() {
  const { userID, updateSpotifyUserID } = useUser();
  const [profile, setProfile] = useState(null);

  useEffect(() => {
    if (userID) {
      fetch(`${config.backendUrl}/spotify/current-profile?userID=${userID}`)
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
    return <Loading />;
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
