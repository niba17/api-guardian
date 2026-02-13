/** @type {import('tailwindcss').Config} */
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}"],
  theme: {
    extend: {
      colors: {
        guardian: {
          dark: "#020617",
          card: "#0f172a",
          primary: "#10b981",
          danger: "#ef4444",
          accent: "#3b82f6",
        },
      },
    },
  },
  plugins: [],
};
