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
      colors: {
        background: "#004643",
        headline: "#fffffe",
        paragraph: "#abd1c6",
        button: "#f9bc60",
        buttonText: "#001e1d",
        stroke: "#001e1d",
        main: "#e8e4e6",
        highlight: "#f9bc60",
        secondary: "#abd1c6",
        tertiary: "#e16162",
        
      }
    },
  },
  plugins: [],
}
