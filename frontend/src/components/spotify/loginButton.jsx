export const LoginButton = () => {
  // TODO: Refactor and move this function to ../../utils/auth-utils.js
  const handleLogin = () => {
    fetch("http://localhost:8080/auth/spotify/login")
      .then((response) => response.json())
      .then((data) => (window.location = data.authURL));
  };

  return (
    <button
      className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 my-5 rounded-full border-2 border-black"
      onClick={handleLogin}
    >
      Connect Spotify Account
    </button>
  );
};
