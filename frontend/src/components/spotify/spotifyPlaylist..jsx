export default function SpotifyPlaylist({ playlist }) {
  return (
    <div id="spotify-playlist">
      <h3>{playlist.name}</h3>
      <p>{playlist.description || "No description available."}</p>
      <p>Tracks: {playlist.tracks.total}</p>
      <p>Owner: {playlist.owner.display_name}</p>
      {playlist.images[0] && (
        <img
          src={playlist.images[0].url}
          alt={playlist.name}
          style={{ height: 100, width: 100 }}
        />
      )}
      <a
        href={playlist.external_urls.spotify}
        target="_blank"
        rel="noopener noreferrer"
      >
        Open in Spotify
      </a>
    </div>
  );
}
