import { FaGithub, FaLinkedin } from "react-icons/fa";

export default function Footer() {
  return (
    <footer className="bg-gray-600 w-full fixed bottom-0 border-t-2 border-red-500 text-center p-2">
      <div className="text-white flex justify-center space-x-2 lg:space-x-4 mt-2">
        <p className="text-xs text-white">
          © {new Date().getFullYear()} Robert Lieblang
        </p>

        <a
          className="text-xs text-white"
          target="_blank"
          rel="noreferrer noopener"
          href="/privacy"
        >
          Privacy Policy
        </a>
        <div className="flex space-x-2 lg:scale-125 lg:ml-4">
          <a
            href="https://github.com/roblieblang"
            target="_blank"
            rel="noopener noreferrer"
          >
            <FaGithub size={15} />
          </a>
          <a
            href="https://linkedin.com/in/roblieblang"
            target="_blank"
            rel="noopener noreferrer"
          >
            <FaLinkedin size={15} />
          </a>
        </div>
      </div>
    </footer>
  );
}
