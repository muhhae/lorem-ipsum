/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ"
  ],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["light", "dark"],
  },
  plugins: [require("daisyui")],
}

