import { Link } from "react-router-dom";

export default function LinkButton({ to, text, onClick, state, className }) {
  const CLASSNAME =
    "hover:bg-white hover:text-green-500 transition text-sm font-bold rounded bg-customSecondary py-1 px-2";

  // Check if 'to' prop starts with "http" to determine if it's an external link
  const isExternal = typeof to === "string" && to.startsWith("http");

  if (isExternal) {
    // Render as an <a> tag for external links
    return (
      <a
        href={to}
        className={CLASSNAME}
        onClick={onClick}
        target="_blank"
        rel="noopener noreferrer"
      >
        {text}
      </a>
    );
  } else if (to) {
    // Render as a Link for internal navigation
    const linkProps =
      typeof to === "object" ? { to: to.pathname, state: to.state } : { to };
    return (
      <Link {...linkProps} className={CLASSNAME} onClick={onClick}>
        {text}
      </Link>
    );
  } else {
    // Render as a button if no 'to' prop is provided
    return (
      <button className={className} onClick={onClick}>
        {text}
      </button>
    );
  }
}
