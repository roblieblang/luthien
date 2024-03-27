import { FaGithub, FaLinkedin } from "react-icons/fa";

export default function Footer() {
  return (
    <footer className="bg-gray-600 w-full fixed bottom-0 border-t-2 border-red-500 text-center p-2">
      <p className="text-sm">© {new Date().getFullYear()} Robert Lieblang</p>
      <div className="flex justify-center space-x-4 mt-2">
        <a
          href="https://github.com/roblieblang"
          target="_blank"
          rel="noopener noreferrer"
        >
          <FaGithub size={25} />
        </a>
        <a
          href="https://linkedin.com/in/roblieblang"
          target="_blank"
          rel="noopener noreferrer"
        >
          <FaLinkedin size={25} />
        </a>
      </div>
    </footer>
  );
}
