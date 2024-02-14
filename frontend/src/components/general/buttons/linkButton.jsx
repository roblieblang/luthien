import { Link } from "react-router-dom"; // for internal navigation

export default function Linkbutton({ to, text, onClick }) {
  const CLASSNAME =
    "bg-blue-500 hover:bg-blue-700 text-white font-bold py-1 px-2 rounded";

  // Render different component for external links
  if (to.includes("https")) {
    return (
      <a href={to} target="_blank" rel="noopener noreferrer">
        <button className={CLASSNAME} onClick={onClick}>
          {text}
        </button>
      </a>
    );
  }
  return (
    <Link to={to}>
      <button className={CLASSNAME} onClick={onClick}>
        {text}
      </button>
    </Link>
  );
}
