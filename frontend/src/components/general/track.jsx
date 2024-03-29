import he from "he";

export function Track({ track, source, formattedLink }) {
  return (
    <div className="bg-customBG rounded border-customParagraph border-solid border-2 p-2 my-0.5 flex lg:w-1/2 w-11/12">
      <div className="flex-none">
        {(track.thumbnailUrl || track.thumbnail) && (
          <img
            src={track.thumbnailUrl || track.thumbnail}
            alt={he.decode(track.title)}
            className="lg:h-28 lg:w-28 h-14 w-14 object-cover mr-2 border-2 rounded"
          />
        )}
      </div>
      <div className="flex-1 text-center flex flex-col justify-center">
        <a
          href={track.link || formattedLink}
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block"
        >
          <h3 className="lg:text-2xl text-sm font-bold text-slate-200 hover:text-blue-500">
            {he.decode(track.title)}
          </h3>
        </a>
        <div className="lg:text-lg text-xs">
          <span>{source.charAt(0).toUpperCase() + source.slice(1)}</span>
          <span className="mx-2">•</span>
          {source === "spotify" ? (
            <>
              <span>{he.decode(track.artist)}</span>
              <span className="mx-2">•</span>
              <span>{he.decode(track.album)}</span>{" "}
            </>
          ) : (
            <span>{he.decode(track.channelTitle)}</span>
          )}
        </div>
      </div>
    </div>
  );
}

export function ModalTrack({ track, destination, formattedLink, isHit }) {
  const title = track.title || track.songTitle;
  const artist = track.artistName || track.artist;

  if (isHit) {
    return (
      <div
        className={`bg-customBG rounded border-2 p-2 my-0.5 flex w-full ${
          isHit ? "border-green-600" : "border-red-600"
        }`}
      >
        <div className="flex-none">
          {track.thumbnail && (
            <img
              src={track.thumbnail}
              alt={he.decode(title)}
              className={`lg:h-28 lg:w-28 h-14 w-14 object-cover mr-2 border-2 rounded ${
                isHit ? "border-green-700" : "border-red-800"
              }`}
            />
          )}
        </div>
        <div className="flex-1 text-center flex flex-col max-w-full">
          <a
            href={formattedLink}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-block"
          >
            <h3 className="text-xs sm:text-sm md:text-lg lg:text-2xl font-bold text-slate-200 hover:text-blue-500 break-all">
              {he.decode(title)}
            </h3>
          </a>
          <div className="text-xs lg:text-lg break-words">
            {destination === "spotify" ? (
              <>
                <span>{he.decode(artist)}</span>
                <span className="mx-1">•</span>
                <span>{he.decode(track.album)}</span>
              </>
            ) : (
              <span>{he.decode(track.channelTitle)}</span>
            )}
          </div>
        </div>
      </div>
    );
  } else {
    return (
      <div
        className={`bg-customBG rounded border-2 p-2 my-0.5 flex w-full ${
          isHit ? "border-green-600" : "border-red-600"
        }`}
      >
        <h3 className="text-xs sm:text-sm md:text-lg lg:text-2xl font-bold text-slate-200 hover:text-blue-500 break-all">
          {he.decode(title)}
        </h3>
      </div>
    );
  }
}
