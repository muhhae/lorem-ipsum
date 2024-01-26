/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./**/*.templ"
  ],
  theme: {
    extend: {
      backgroundImage: theme => ({
        'gradient-to-l': `linear-gradient(to left, ${theme('colors.primary.DEFAULT')} 0%, transparent 100%)`
      })
    }
  },
  daisyui: {
    themes: ["nord", "sunset"],
  },
  plugins: [require("daisyui")],
}

