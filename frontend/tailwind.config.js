/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    screens: {
      "sm": "440px",
      "md": "547px",
      "lg": "768px",
      "xl": "1024px",
      "2xl": "1680px",
    },
    extend: {
      fontFamily: {
        // "custom": ["Playfair Display", "serif"],
      },
      colors: {
        customBG: "#0a1c1d",
        customHeadline: "#202d27",
        customButton: "#95b2f1",
        customStroke: "#e2714a",
        customParagraph: "#94a1b2",
        customButtonText: "#fffffe",
        customPrimary: "#fffffe",
        customHighlight: "#e2714a",
        customSecondary: "#72757e",
        customTertiary: "#2cb67d",
      }
    },
  },
  plugins: [],
}
