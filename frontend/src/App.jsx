import { Route, Routes } from "react-router-dom";
import Callback from "./pages/callback";
import Home from "./pages/home";
import Music from "./pages/music";
import Profile from "./pages/profile";

export default function App() {
  return (
    <>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/profile" element={<Profile />} />
        <Route path="/callback" element={<Callback />} />
        <Route path="/music" element={<Music />} />
      </Routes>
    </>
  );
}
