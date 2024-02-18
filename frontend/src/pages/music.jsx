import LinkButton from "../components/general/buttons/linkButton";
import BasicHeading from "../components/general/headings/basicHeading";
import SpotifyUserProfile from "../components/spotify/spotifyUserProfile";
import YouTubeUserProfile from "../components/youtube/youtubeUserProfile";

export default function Music() {
  return (
    <div className="mb-2">
      <BasicHeading text="Music Page" />
      <SpotifyUserProfile />
      <YouTubeUserProfile />
      <div className="my-5">
        <LinkButton to="/" text="Back" />
      </div>
    </div>
  );
}
