import { Route, Routes } from "react-router-dom";
import ProtectedRoute from "./components/auth0/protectedRoute";
import Conversion from "./pages/conversion";
import GoogleCallback from "./pages/googleCallback";
import Home from "./pages/home";
import Music from "./pages/music";
import Profile from "./pages/profile";
import SpotifyCallback from "./pages/spotifyCallback";

export default function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/spotify/callback" element={<SpotifyCallback />} />
        <Route path="/google/callback" element={<GoogleCallback />} />
        <Route
          path="/music"
          element={
            <ProtectedRoute>
              <Music />
            </ProtectedRoute>
          }
        />
        <Route
          path="/conversion"
          element={
            <ProtectedRoute>
              <Conversion />
            </ProtectedRoute>
          }
        />
      </Routes>
    </>
  );
}
