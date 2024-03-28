import { useAuth0 } from "@auth0/auth0-react";
import Loading from "../general/modals/loading";
import UserPlaylists from "../general/userPlaylists";

export default function YouTubeUserProfile() {
  const { user } = useAuth0();

  if (!user) {
    return <Loading />;
  }

  return (
    <div>
      <h1 className="font-bold text-lg">{user.name}&apos;s YouTube Account</h1>
      <UserPlaylists serviceType={"youtube"} />
    </div>
  );
}
