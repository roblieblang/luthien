import { useEffect } from "react";
import { useUser } from "../contexts/userContext";

export default function GoogleCallback() {
  // const { userID } = useUser();
  const userID = sessionStorage.getItem("userID");

  useEffect(() => {
    if (window.location.pathname === "/google/callback") {
      console.log(`\nApp UID: ${userID}`);
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get("code");
      const sessionID = sessionStorage.getItem("sessionID");

      fetch("http://localhost:8080/auth/google/callback", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ code, userID, sessionID }),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data.redirectURL) {
            window.location.href = data.redirectURL;
          }
        })
        .catch((error) => {
          console.error("Error during Google callback processing:", error);
        });
    }
  });

  return <h1>Loading...</h1>;
}
