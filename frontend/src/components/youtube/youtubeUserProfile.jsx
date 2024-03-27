import { useAuth0 } from "@auth0/auth0-react";
import { Bars } from "react-loader-spinner";
import UserPlaylists from "../general/userPlaylists";

export default function YouTubeUserProfile() {
  const { user } = useAuth0();

  if (!user) {
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
      <h1 className="font-bold text-lg">{user.name}&apos;s YouTube Account</h1>
      <UserPlaylists serviceType={"youtube"} />
    </div>
  );
}
