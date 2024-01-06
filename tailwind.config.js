/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ"
  ],
  theme: {
    extend: {},
  },
  daisyui: {
    themes: ["nord", "sunset"],
  },
  plugins: [require("daisyui")],
}

